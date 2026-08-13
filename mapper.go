package fairy

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/kirinyoku/enkanetwork-go/client/zzz"
	"github.com/kirinyoku/fairy/store"
)

// profileMapper is responsible for translating raw EnkaNetwork API data into
// clean fairy domain models using a specific language and metadata store.
type profileMapper struct {
	store store.MetadataStore
	lang  Language
}

// newMapper creates a new mapper instance.
func newMapper(s store.MetadataStore, lang Language) *profileMapper {
	return &profileMapper{
		store: s,
		lang:  lang,
	}
}

// ToProfile converts a raw zzz.Profile into an enriched Profile.
func (m *profileMapper) ToProfile(raw *zzz.Profile) (*Profile, error) {
	if raw == nil {
		return nil, fmt.Errorf("raw profile is nil")
	}

	p := &Profile{
		Region: mapRegion(raw.Region),
		Agents: make([]Agent, 0),
	}

	if raw.PlayerInfo.SocialDetail != nil {
		if pd := raw.PlayerInfo.SocialDetail.ProfileDetail; pd != nil {
			p.Nickname = pd.Nickname
			p.InterknotLevel = pd.Level
			p.UID = strconv.FormatInt(pd.UID, 10)

			p.Title = m.mapTitle(pd)
			p.Avatar = m.mapAvatar(pd.ProfileID, pd.AvatarID)
			p.Namecard = m.mapNamecard(pd.CallingCardID)
		}

		p.Badges = make([]Badge, 0, len(raw.PlayerInfo.SocialDetail.MedalList))
		for _, rawMedal := range raw.PlayerInfo.SocialDetail.MedalList {
			if badge := m.mapBadge(rawMedal); badge != nil {
				p.Badges = append(p.Badges, *badge)
			}
		}
	}

	if raw.PlayerInfo.ShowcaseDetail != nil {
		p.Agents = make([]Agent, 0, len(raw.PlayerInfo.ShowcaseDetail.AvatarList))
		for _, rawAvatar := range raw.PlayerInfo.ShowcaseDetail.AvatarList {
			agent := m.ToAgent(&rawAvatar)
			if agent != nil {
				p.Agents = append(p.Agents, *agent)
			}
		}
	}

	return p, nil
}

func (m *profileMapper) mapTitle(pd *zzz.ProfileDetail) *Title {
	if pd == nil {
		return nil
	}

	titleID := pd.Title
	var args []string

	if pd.TitleInfo != nil {
		if pd.TitleInfo.FullTitle != 0 {
			titleID = pd.TitleInfo.FullTitle
		} else if pd.TitleInfo.Title != 0 {
			titleID = pd.TitleInfo.Title
		}
		for _, arg := range pd.TitleInfo.Args {
			args = append(args, fmt.Sprintf("%v", arg))
		}
	}

	if titleID == 0 {
		return nil
	}

	meta, ok := m.store.TitleMeta(titleID)
	if !ok && pd.TitleInfo != nil && pd.TitleInfo.Title != 0 && pd.TitleInfo.Title != titleID {
		meta, ok = m.store.TitleMeta(pd.TitleInfo.Title)
	}

	if !ok {
		return &Title{ID: titleID}
	}

	text := m.store.Localize(meta.TitleText, string(m.lang))
	for i, arg := range args {
		text = strings.ReplaceAll(text, fmt.Sprintf("{%d}", i), arg)
	}

	return &Title{
		ID:             titleID,
		Text:           text,
		PrimaryColor:   meta.ColorA,
		SecondaryColor: meta.ColorB,
	}
}

func (m *profileMapper) mapAvatar(profileID, avatarID int) *Avatar {
	if meta, ok := m.store.PfpMeta(profileID); ok && meta.Icon != "" {
		return &Avatar{
			ID:  profileID,
			URL: buildEnkaURL(meta.Icon),
		}
	}
	if meta, ok := m.store.PfpMeta(avatarID); ok && meta.Icon != "" {
		return &Avatar{
			ID:  avatarID,
			URL: buildEnkaURL(meta.Icon),
		}
	}
	if avatarMeta, ok := m.store.AvatarMeta(avatarID); ok && avatarMeta.CircleIcon != "" {
		return &Avatar{
			ID:  avatarID,
			URL: buildEnkaURL(avatarMeta.CircleIcon),
		}
	}
	if skinMeta, ok := m.store.SkinMeta(avatarID); ok && skinMeta.Icon != "" {
		return &Avatar{
			ID:  avatarID,
			URL: buildEnkaURL(skinMeta.Icon),
		}
	}
	if defaultSkin, ok := m.store.DefaultSkinMeta(avatarID); ok && defaultSkin.Icon != "" {
		return &Avatar{
			ID:  avatarID,
			URL: buildEnkaURL(defaultSkin.Icon),
		}
	}
	return &Avatar{ID: profileID}
}

func (m *profileMapper) mapNamecard(id int) *Namecard {
	meta, ok := m.store.NamecardMeta(id)
	if !ok {
		return &Namecard{ID: id}
	}
	return &Namecard{
		ID:  id,
		URL: buildEnkaURL(meta.Icon),
	}
}

func (m *profileMapper) mapBadge(raw zzz.Medal) *Badge {
	meta, ok := m.store.MedalMeta(raw.MedalIcon)
	if !ok {
		return &Badge{ID: raw.MedalIcon, Value: raw.Value}
	}
	return &Badge{
		ID:      raw.MedalIcon,
		Title:   m.store.Localize(meta.Name, string(m.lang)),
		Value:   raw.Value,
		IconURL: buildEnkaURL(meta.Icon),
	}
}

// ToAgent converts a raw zzz.AvatarData into an enriched Agent.
// The mapping pipeline performs the following steps:
// 1. Basic properties mapping (Level, Promotion, etc.)
// 2. Metadata lookup (Name, Attribute, Specialty) with fallback to "Unknown"
// 3. Skin resolution (equipped vs default)
// 4. Equipment mapping (W-Engine and Drive Discs)
// 5. Set bonus calculation based on equipped Discs
// 6. Skill level mapping (handling mismatches between raw levels and template skills)
// 7. Final stat calculation.
func (m *profileMapper) ToAgent(raw *zzz.AvatarData) *Agent {
	if raw == nil {
		return nil
	}

	agent := &Agent{
		ID:                   raw.ID,
		Level:                raw.Level,
		Promotion:            raw.PromotionLevel,
		MindscapeCinema:      raw.TalentLevel,
		CoreSkillEnhancement: raw.CoreSkillEnhancement,
	}

	meta, ok := m.store.AvatarMeta(raw.ID)
	if ok {
		agent.Name = m.store.Localize(meta.Name, string(m.lang))
		agent.Rarity = mapRarity(meta.Rarity)
		agent.SplashArtURL = buildEnkaURL(meta.Image)

		if len(meta.ElementTypes) > 0 {
			agent.Attribute = mapRawToAttribute(meta.ElementTypes[0])
			agent.AttributeName = m.store.Localize("ElementType_"+meta.ElementTypes[0], string(m.lang))
		}

		agent.Specialty = Specialty(meta.ProfessionType)
		agent.SpecialtyName = m.store.Localize(mapSpecialtyLocKey(meta.ProfessionType), string(m.lang))
	} else {
		agent.Name = fmt.Sprintf("Unknown Agent (%d)", raw.ID)
	}

	var skinMeta store.SkinMeta
	var foundSkin bool
	if raw.SkinID != 0 {
		skinMeta, foundSkin = m.store.SkinMeta(raw.SkinID)
	} else {
		skinMeta, foundSkin = m.store.DefaultSkinMeta(raw.ID)
	}

	if foundSkin {
		agent.Skin = &Skin{
			ID:           skinMeta.Id,
			Name:         m.store.Localize(skinMeta.NameKey, string(m.lang)),
			Description:  m.store.Localize(skinMeta.DescKey, string(m.lang)),
			SplashArtURL: buildEnkaURL(skinMeta.Icon),
		}
	}

	if raw.Weapon != nil {
		agent.WEngine = m.mapWEngine(raw.Weapon)
	}

	agent.DriveDiscs = make([]DriveDisc, 0, len(raw.EquippedList))
	for _, eq := range raw.EquippedList {
		if eq.Equipment != nil {
			disc := m.mapDriveDisc(eq.Equipment, eq.Slot)
			if disc != nil {
				agent.DriveDiscs = append(agent.DriveDiscs, *disc)
			}
		}
	}
	agent.ActiveSetBonuses = m.detectSetBonuses(agent.DriveDiscs)

	if skillsMeta, ok := m.store.AvatarSkillsMeta(raw.ID); ok {
		// Skill mapping logic is complex because the API returns only 5 skill levels
		// (0: Basic, 1: Dodge, 2: Assist, 3: Special, 4: Chain) in the raw `SkillLevelList`,
		// but the game's datamined template `AvatarSkillDesTemplateTb` contains 12+ active skills.
		// `groupMap` is used to map a specific template skill index back to the appropriate skill group index
		// in the API's `SkillLevelList`.
		// The typical order of 12 active skills is:
		// 0,1 (Basic), 2,3 (Special), 4,5,6 (Dodge), 7,8 (Chain), 9,10,11 (Assist)
		groupMap := map[int]int{
			0: 0, 1: 0,
			2: 3, 3: 3,
			4: 1, 5: 1, 6: 1,
			7: 4, 8: 4,
			9: 2, 10: 2, 11: 2,
		}

		// Pre-process skill levels from raw data
		levels := make(map[int]int)
		for _, sl := range raw.SkillLevelList {
			levels[sl.Index] = sl.Level
		}

		agent.Skills = make([]Skill, 0, len(skillsMeta))
		for i, sm := range skillsMeta {
			lvl := 0
			if group, found := groupMap[i]; found {
				lvl = levels[group]
			} else if i >= 12 {
				// Passives
				lvl = raw.CoreSkillEnhancement
			}

			agent.Skills = append(agent.Skills, Skill{
				Level:       lvl,
				Name:        m.store.Localize(sm.NameKey, string(m.lang)),
				Description: m.store.Localize(sm.DescKey, string(m.lang)),
			})
		}
	}

	if msMeta, ok := m.store.AvatarMindscapesMeta(raw.ID); ok {
		mindscapes := make([]MindscapeNode, 0, len(msMeta))
		for _, mn := range msMeta {
			mindscapes = append(mindscapes, MindscapeNode{
				Rank:        mn.Rank,
				Name:        m.store.Localize(mn.TitleKey, string(m.lang)),
				Description: m.store.Localize(mn.DescKey, string(m.lang)),
				Unlocked:    raw.TalentLevel >= mn.Rank,
			})
		}
		agent.Mindscapes = mindscapes
	}

	if pvMeta, ok := m.store.AvatarPotentialVisionsMeta(raw.ID); ok {
		isUnlocked := false
		if raw.IsUpgradeUnlocked != nil {
			isUnlocked = *raw.IsUpgradeUnlocked
		} else if raw.UpgradeID > 0 {
			isUnlocked = true
		}

		pvNodes := make([]PotentialVisionNode, 0, len(pvMeta))
		for _, nodeMeta := range pvMeta {
			isActive := isUnlocked && (raw.UpgradeID == nodeMeta.ID || (raw.UpgradeID > 0 && raw.UpgradeID >= nodeMeta.ID))
			pvNodes = append(pvNodes, PotentialVisionNode{
				ID:          nodeMeta.ID,
				Level:       nodeMeta.Level,
				LevelName:   m.store.Localize(nodeMeta.LevelNameKey, string(m.lang)),
				Title:       m.store.Localize(nodeMeta.TitleKey, string(m.lang)),
				Description: m.store.Localize(nodeMeta.DescKey, string(m.lang)),
				IsActive:    isActive,
			})
		}

		agent.PotentialVision = &PotentialVision{
			IsUnlocked: isUnlocked,
			CurrentID:  raw.UpgradeID,
			Nodes:      pvNodes,
		}
	}

	// Calculate final stats
	calculateAgentStats(agent, m.store)

	return agent
}

// mapWEngine converts raw weapon data into a WEngine domain object.
// `PassiveDescKeys` from WeaponMeta contains a slice of description keys
// corresponding to the weapon's phase/refinement level (0-4). We use the (1-indexed)
// Modification level - 1 to look up the correct description text.
func (m *profileMapper) mapWEngine(raw *zzz.Weapon) *WEngine {
	w := &WEngine{
		ID:           raw.ID,
		UID:          strconv.Itoa(raw.UID),
		Level:        raw.Level,
		Phase:        raw.BreakLevel,
		Modification: raw.UpgradeLevel,
	}

	meta, ok := m.store.WeaponMeta(raw.ID)
	if ok {
		w.Name = m.store.Localize(meta.ItemName, string(m.lang))
		w.Rarity = mapRarity(meta.Rarity)
		w.Specialty = Specialty(meta.ProfessionType)
		w.SpecialtyName = m.store.Localize(mapSpecialtyLocKey(meta.ProfessionType), string(m.lang))
		w.IconURL = buildEnkaURL(meta.ImagePath)

		wPhase := raw.UpgradeLevel - 1
		if wPhase < 0 {
			wPhase = 0
		}
		w.MainStat = m.mapStat(meta.MainStat.PropertyID, calcWEngineMainStat(m.store, meta, raw.Level, raw.BreakLevel), 0)
		w.SecondaryStat = m.mapStat(meta.SecondaryStat.PropertyID, calcWEngineSecondaryStat(m.store, meta, raw.Level, wPhase), 0)

		if wPhase < len(meta.PassiveDescKeys) {
			w.PassiveDescription = m.store.Localize(meta.PassiveDescKeys[wPhase], string(m.lang))
		}
	}

	return w
}

// mapDriveDisc converts raw equipment data into a DriveDisc domain object.
// The main stat scaling works differently here. The raw API provides
// only the base value of the main stat at Level 0. We must look up the correct
// level modifier from the EquipmentLevelTemplate and apply it (LevelMod / 10000)
// to calculate the actual main stat value at the current disc level.
func (m *profileMapper) mapDriveDisc(raw *zzz.Equipment, slot int) *DriveDisc {
	d := &DriveDisc{
		ID:    raw.ID,
		UID:   strconv.Itoa(raw.UID),
		Slot:  slot,
		Level: raw.Level,
	}

	meta, ok := m.store.EquipmentMeta(raw.ID)
	if ok {
		d.Rarity = mapRarity(meta.Rarity)
		d.Set = Set{ID: meta.SuitID}

		if suitMeta, suitOk := m.store.EquipmentSuitMeta(meta.SuitID); suitOk {
			d.Set.Name = m.store.Localize(suitMeta.Name, string(m.lang))
			d.IconPath = buildEnkaURL(suitMeta.Icon)
		}
	}

	if len(raw.MainPropertyList) > 0 {
		prop := raw.MainPropertyList[0]

		levelMod := 0
		if meta, ok := m.store.EquipmentMeta(raw.ID); ok {
			if lvlTpl, ok := m.store.EquipmentLevelTemplate(meta.Rarity, raw.Level); ok {
				levelMod = lvlTpl.MainStat
			}
		}

		scaledVal := float64(prop.PropertyValue) * (1.0 + float64(levelMod)/10000.0)
		d.MainStat = m.mapStat(prop.PropertyID, int(scaledVal), 0)
	}

	d.SubStats = make([]StatValue, 0, len(raw.RandomPropertyList))
	for _, prop := range raw.RandomPropertyList {
		rolls := prop.PropertyLevel
		if rolls == 0 {
			rolls = 1
		}
		d.SubStats = append(d.SubStats, m.mapStat(prop.PropertyID, prop.PropertyValue*rolls, rolls))
	}

	return d
}

// mapStat creates a StatValue from a property ID, raw value, and roll count.
// It detects whether a stat is a percentage by inspecting the `Format` string
// in the PropertyMeta (e.g., "{0:0.#}%"). If it is a percentage, the raw integer
// value is divided by 10000.0, as per the game's internal representation convention.
func (m *profileMapper) mapStat(propID int, rawValue int, rolls int) StatValue {
	pid := PropertyID(propID)
	sv := StatValue{
		PropertyID: pid,
		Rolls:      rolls,
		Value:      float64(rawValue),
		IconURL:    pid.IconURL(),
	}

	meta, ok := m.store.PropertyMeta(propID)
	if ok {
		sv.Name = m.store.Localize(meta.Name, string(m.lang))
		sv.IsPercent = strings.Contains(meta.Format, "%")
		if sv.IsPercent {
			sv.Value = float64(rawValue) / 10000.0
		}
	} else {
		sv.Name = fmt.Sprintf("Stat_%d", propID)
	}

	return sv
}

// detectSetBonuses returns the active set bonuses from equipped drive discs.
// This reports the highest achieved set bonus (2pc or 4pc) in the UI text,
// but both 2pc and 4pc property bonuses are accumulated behind the scenes.
func (m *profileMapper) detectSetBonuses(discs []DriveDisc) []DriveDiscSetBonus {
	counts := make(map[int]int)
	names := make(map[int]string)

	for _, d := range discs {
		if d.Set.ID > 0 {
			counts[d.Set.ID]++
			names[d.Set.ID] = d.Set.Name
		}
	}

	var bonuses []DriveDiscSetBonus
	for setID, count := range counts {
		if count < 2 {
			continue
		}

		var setMeta store.EquipmentSuitMeta
		if meta, ok := m.store.EquipmentSuitMeta(setID); ok {
			setMeta = meta
		}

		var effects []SetEffect
		if setMeta.Set2DescKey != "" {
			effects = append(effects, SetEffect{
				PieceCount:  2,
				Description: m.store.Localize(setMeta.Set2DescKey, string(m.lang)),
				IsActive:    count >= 2,
			})
		}
		if setMeta.Set4DescKey != "" {
			effects = append(effects, SetEffect{
				PieceCount:  4,
				Description: m.store.Localize(setMeta.Set4DescKey, string(m.lang)),
				IsActive:    count >= 4,
			})
		}

		bonuses = append(bonuses, DriveDiscSetBonus{
			Set:     Set{ID: setID, Name: names[setID]},
			Count:   count,
			Effects: effects,
		})
	}

	return bonuses
}

// mapSpecialtyLocKey maps the raw profession type string to the locs.json key.
func mapSpecialtyLocKey(raw string) string {
	switch raw {
	case "Attack":
		return "ProfessionName_PowerfulAttack"
	case "Stun":
		return "ProfessionName_BreakStun"
	case "Anomaly":
		return "ProfessionName_ElementAbnormal"
	case "Support":
		return "ProfessionName_Support"
	case "Defense":
		return "ProfessionName_Defence"
	case "Rupture":
		return "ProfessionName_Rupture"
	default:
		return "ProfessionName_" + raw
	}
}

// mapRawToAttribute maps the internal element type string to the domain Attribute enum.
func mapRawToAttribute(raw string) Attribute {
	switch raw {
	case "Physics":
		return AttributePhysical
	case "Elec":
		return AttributeElectric
	case "Fire":
		return AttributeFire
	case "Ice":
		return AttributeIce
	case "Ether":
		return AttributeEther
	case "Wind":
		return AttributeWind
	case "FireFrost":
		return AttributeFrost
	case "ZhenZhenAssault":
		return AttributeHonedEdge
	case "AuricEther":
		return AttributeAuricInk
	case "Lumen":
		return AttributeLumiflux
	default:
		return Attribute(raw)
	}
}

// mapRarity maps the internal integer rarity tier to the string enum.
func mapRarity(raw int) Rarity {
	switch raw {
	case 4:
		return RarityS
	case 3:
		return RarityA
	case 2:
		return RarityB
	default:
		return Rarity("Unknown")
	}
}

// mapRegion maps the raw region string from EnkaNetwork to the domain Region enum.
func mapRegion(raw string) Region {
	rawLower := strings.ToLower(raw)
	switch rawLower {
	case "eu", "europe", "prod_gf_eu":
		return RegionEU
	case "na", "america", "usa", "prod_gf_us":
		return RegionNA
	case "asia", "prod_gf_jp":
		return RegionAsia
	case "twhkmo", "tw", "hk", "mo", "prod_gf_sg":
		return RegionTWHKMO
	default:
		// Fallback
		return Region(raw)
	}
}

// buildEnkaURL constructs a full URL to the Enka.network ZZZ image asset.
func buildEnkaURL(path string) string {
	if path == "" {
		return ""
	}
	if strings.HasPrefix(path, "http") {
		return path
	}
	if !strings.HasPrefix(path, "/") {
		if !strings.HasPrefix(path, "ui/") {
			path = "/ui/zzz/" + path
		} else {
			path = "/" + path
		}
	}
	if !strings.HasSuffix(path, ".png") && !strings.HasSuffix(path, ".jpg") && !strings.HasSuffix(path, ".webp") {
		path = path + ".png"
	}
	return "https://enka.network" + path
}
