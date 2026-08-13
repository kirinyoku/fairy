package store

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"strconv"
	"sync"

	"github.com/kirinyoku/fairy/internal/assets"
)

// EmbeddedStore implements MetadataStore using the embedded JSON files.
// It relies on the go:embed feature to package large JSON datamined files
// directly into the binary. This ensures that the library requires no external
// file dependencies at runtime and loads instantaneously into memory.
type EmbeddedStore struct {
	avatars              map[int]AvatarMeta
	weapons              map[int]WeaponMeta
	equipments           map[int]EquipmentMeta
	equipmentSuits       map[int]EquipmentSuitMeta
	properties           map[int]PropertyMeta
	locs                 map[string]map[string]string
	medals               map[int]MedalMeta
	titles               map[int]TitleMeta
	namecards            map[int]NamecardMeta
	pfps                 map[int]PfpMeta
	skins                map[int]SkinMeta
	defaultSkins         map[int]SkinMeta
	skills               map[int][]SkillMeta
	skillTemplates       map[int]SkillTemplateMeta
	mindscapes           map[int][]MindscapeMeta
	potentialVisions     map[int][]PotentialVisionMeta
	weaponLevelTemplates map[templateKey]WeaponLevelTemplate
	weaponStarTemplates  map[templateKey]WeaponStarTemplate
	equipLevelTemplates  map[templateKey]EquipmentLevelTemplate
}

// Ensure EmbeddedStore implements MetadataStore.
var _ MetadataStore = (*EmbeddedStore)(nil)

var (
	defaultStore     *EmbeddedStore
	defaultStoreOnce sync.Once
	defaultStoreErr  error
)

// Default returns a lazily initialized singleton EmbeddedStore.
// Uses the singleton pattern with sync.Once because parsing the embedded
// JSON files takes some CPU and memory. Since the store is completely read-only
// after initialization, sharing a single instance across the entire application is safe.
func Default() (*EmbeddedStore, error) {
	defaultStoreOnce.Do(func() {
		subFS, err := fs.Sub(assets.DataFS, "data")
		if err != nil {
			defaultStoreErr = fmt.Errorf("failed to sub embedded FS: %w", err)
			return
		}
		defaultStore, defaultStoreErr = parseStore(subFS)
	})
	return defaultStore, defaultStoreErr
}

// readJSON is a generic helper that abstracts the boilerplate of reading and unmarshaling
// JSON from the embedded filesystem into strongly typed structs.
func readJSON[T any](fsys fs.FS, path string) (T, error) {
	var result T
	data, err := fs.ReadFile(fsys, path)
	if err != nil {
		return result, fmt.Errorf("failed to read %s: %w", path, err)
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return result, fmt.Errorf("failed to unmarshal %s: %w", path, err)
	}
	return result, nil
}

// parseStore builds the EmbeddedStore by loading all embedded JSON files.
// It loads all files sequentially. Currently, it returns an error if any
// file fails to load to ensure data consistency, as missing a file
// would result in silently incorrect stat calculations.
func parseStore(fsys fs.FS) (*EmbeddedStore, error) {
	s := &EmbeddedStore{
		avatars:              make(map[int]AvatarMeta),
		weapons:              make(map[int]WeaponMeta),
		equipments:           make(map[int]EquipmentMeta),
		equipmentSuits:       make(map[int]EquipmentSuitMeta),
		properties:           make(map[int]PropertyMeta),
		medals:               make(map[int]MedalMeta),
		titles:               make(map[int]TitleMeta),
		namecards:            make(map[int]NamecardMeta),
		pfps:                 make(map[int]PfpMeta),
		skins:                make(map[int]SkinMeta),
		defaultSkins:         make(map[int]SkinMeta),
		skills:               make(map[int][]SkillMeta),
		skillTemplates:       make(map[int]SkillTemplateMeta),
		mindscapes:           make(map[int][]MindscapeMeta),
		potentialVisions:     make(map[int][]PotentialVisionMeta),
		weaponLevelTemplates: make(map[templateKey]WeaponLevelTemplate),
		weaponStarTemplates:  make(map[templateKey]WeaponStarTemplate),
		equipLevelTemplates:  make(map[templateKey]EquipmentLevelTemplate),
	}

	// parseID is a helper closure for DRYing up string-to-int conversion.
	// It injects contextual error messages to make debugging malformed JSONs easier.
	parseID := func(str string, context string) (int, error) {
		id, err := strconv.Atoi(str)
		if err != nil {
			return 0, fmt.Errorf("invalid %s id %q: %w", context, str, err)
		}
		return id, nil
	}

	// Avatars
	avatarRecs, err := readJSON[map[string]avatarRecord](fsys, "avatars.json")
	if err != nil {
		return nil, err
	}
	for idStr, rec := range avatarRecs {
		id, err := parseID(idStr, "avatar")
		if err != nil {
			return nil, err
		}

		// BaseProps, GrowthProps, PromotionProps, and CoreEnhancementProps all come
		// as string-keyed maps from JSON. Convert them to integer-keyed maps
		// to avoid expensive string conversions during hot-path stat calculations.
		baseProps := make(map[int]int, len(rec.BaseProps))
		for k, v := range rec.BaseProps {
			propID, err := parseID(k, "avatar base prop")
			if err != nil {
				return nil, err
			}
			baseProps[propID] = v
		}

		growthProps := make(map[int]int, len(rec.GrowthProps))
		for k, v := range rec.GrowthProps {
			propID, err := parseID(k, "avatar growth prop")
			if err != nil {
				return nil, err
			}
			growthProps[propID] = v
		}

		promProps := make([]map[int]int, len(rec.PromotionProps))
		for i, pMap := range rec.PromotionProps {
			promProps[i] = make(map[int]int, len(pMap))
			for k, v := range pMap {
				propID, err := parseID(k, "avatar promotion prop")
				if err != nil {
					return nil, err
				}
				promProps[i][propID] = v
			}
		}

		coreProps := make([]map[int]int, len(rec.CoreEnhancementProps))
		for i, pMap := range rec.CoreEnhancementProps {
			coreProps[i] = make(map[int]int, len(pMap))
			for k, v := range pMap {
				propID, err := parseID(k, "avatar core prop")
				if err != nil {
					return nil, err
				}
				coreProps[i][propID] = v
			}
		}

		s.avatars[id] = AvatarMeta{
			Name:                 rec.Name,
			Rarity:               rec.Rarity,
			ProfessionType:       rec.ProfessionType,
			ElementTypes:         rec.ElementTypes,
			BaseProps:            baseProps,
			GrowthProps:          growthProps,
			PromotionProps:       promProps,
			CoreEnhancementProps: coreProps,
			Image:                rec.Image,
			CircleIcon:           rec.CircleIcon,
		}
	}

	// Weapons
	weaponRecs, err := readJSON[map[string]weaponRecord](fsys, "weapons.json")
	if err != nil {
		return nil, err
	}
	for idStr, rec := range weaponRecs {
		id, err := parseID(idStr, "weapon")
		if err != nil {
			return nil, err
		}
		s.weapons[id] = WeaponMeta(rec)
	}

	// Equipments
	eqCont, err := readJSON[equipmentContainer](fsys, "equipments.json")
	if err != nil {
		return nil, err
	}
	for idStr, rec := range eqCont.Items {
		id, err := parseID(idStr, "equipment")
		if err != nil {
			return nil, err
		}
		s.equipments[id] = EquipmentMeta{Rarity: rec.Rarity, SuitID: rec.SuitId}
	}
	for idStr, rec := range eqCont.Suits {
		id, err := parseID(idStr, "equipment suit")
		if err != nil {
			return nil, err
		}
		props := make(map[int]int, len(rec.SetBonusProps))
		for k, v := range rec.SetBonusProps {
			propID, err := parseID(k, "equipment suit prop")
			if err != nil {
				return nil, err
			}
			props[propID] = v
		}
		s.equipmentSuits[id] = EquipmentSuitMeta{
			Icon:          rec.Icon,
			Name:          rec.Name,
			SetBonusProps: props,
			Set2DescKey:   rec.Set2DescKey,
			Set4DescKey:   rec.Set4DescKey,
		}
	}

	// Properties
	propertiesRecs, err := readJSON[map[string]PropertyMeta](fsys, "property.json")
	if err != nil {
		return nil, err
	}
	for idStr, rec := range propertiesRecs {
		id, err := parseID(idStr, "property")
		if err != nil {
			return nil, err
		}
		s.properties[id] = rec
	}

	// Locs
	locs, err := readJSON[map[string]map[string]string](fsys, "locs.json")
	if err != nil {
		return nil, err
	}
	s.locs = locs

	// Titles
	titles, err := readJSON[titleContainer](fsys, "titles.json")
	if err != nil {
		return nil, err
	}
	for idStr, rec := range titles.Titles {
		id, err := parseID(idStr, "title")
		if err != nil {
			return nil, err
		}
		s.titles[id] = rec
	}
	for idStr, rec := range titles.TitleVariants {
		id, err := parseID(idStr, "title variant")
		if err != nil {
			return nil, err
		}
		if rec.ColorA == "" || rec.ColorB == "" {
			parentID := id / 100
			if parent, parentOk := s.titles[parentID]; parentOk {
				if rec.ColorA == "" {
					rec.ColorA = parent.ColorA
				}
				if rec.ColorB == "" {
					rec.ColorB = parent.ColorB
				}
			}
		}
		s.titles[id] = rec
	}

	// Namecards
	namecards, err := readJSON[map[string]NamecardMeta](fsys, "namecards.json")
	if err != nil {
		return nil, err
	}
	for idStr, nc := range namecards {
		id, err := parseID(idStr, "namecard")
		if err != nil {
			return nil, err
		}
		s.namecards[id] = nc
	}

	// Skins
	skins, err := readJSON[map[string]SkinMeta](fsys, "skins.json")
	if err != nil {
		return nil, err
	}
	for idStr, skin := range skins {
		id, err := parseID(idStr, "skin")
		if err != nil {
			return nil, err
		}
		s.skins[id] = skin
		if skin.IsDefault {
			s.defaultSkins[skin.AvatarId] = skin
		}
	}

	// Skills
	skills, err := readJSON[map[string][]SkillMeta](fsys, "skills.json")
	if err != nil {
		return nil, err
	}
	for idStr, skillList := range skills {
		id, err := parseID(idStr, "skill avatar")
		if err != nil {
			return nil, err
		}
		s.skills[id] = skillList
	}

	// Skill Templates
	skillTemplates, err := readJSON[map[string]SkillTemplateMeta](fsys, "skill_templates.json")
	if err == nil {
		for idStr, st := range skillTemplates {
			id, err := parseID(idStr, "skill template")
			if err == nil {
				s.skillTemplates[id] = st
			}
		}
	}

	// Mindscapes
	mindscapes, err := readJSON[map[string][]MindscapeMeta](fsys, "mindscapes.json")
	if err != nil {
		return nil, err
	}
	for idStr, msList := range mindscapes {
		id, err := parseID(idStr, "mindscape avatar")
		if err != nil {
			return nil, err
		}
		s.mindscapes[id] = msList
	}

	// Potential Visions
	potentialVisions, err := readJSON[map[string][]PotentialVisionMeta](fsys, "potential_visions.json")
	if err != nil {
		return nil, err
	}
	for idStr, pvList := range potentialVisions {
		id, err := parseID(idStr, "potential vision avatar")
		if err != nil {
			return nil, err
		}
		s.potentialVisions[id] = pvList
	}

	// Pfps
	pfps, err := readJSON[map[string]PfpMeta](fsys, "pfps.json")
	if err != nil {
		return nil, err
	}
	for idStr, pfp := range pfps {
		id, err := parseID(idStr, "pfp")
		if err != nil {
			return nil, err
		}
		s.pfps[id] = pfp
	}

	// Medals
	medals, err := readJSON[map[string]MedalMeta](fsys, "medals.json")
	if err != nil {
		return nil, err
	}
	for idStr, medal := range medals {
		id, err := parseID(idStr, "medal")
		if err != nil {
			return nil, err
		}
		s.medals[id] = medal
	}

	// Templates
	wlt, err := readJSON[weaponLevelContainer](fsys, "WeaponLevelTemplateTb.json")
	if err != nil {
		return nil, err
	}
	for _, rec := range wlt.List {
		s.weaponLevelTemplates[templateKey{Rarity: rec.Rarity, Level: rec.Level}] = rec
	}

	wst, err := readJSON[weaponStarContainer](fsys, "WeaponStarTemplateTb.json")
	if err != nil {
		return nil, err
	}
	for _, rec := range wst.List {
		s.weaponStarTemplates[templateKey{Rarity: rec.Rarity, Level: rec.BreakLevel}] = rec
	}

	elt, err := readJSON[equipmentLevelContainer](fsys, "EquipmentLevelTemplateTb.json")
	if err != nil {
		return nil, err
	}
	for _, rec := range elt.List {
		s.equipLevelTemplates[templateKey{Rarity: rec.Rarity, Level: rec.Level}] = rec
	}

	return s, nil
}

// Localize translates a text hash into the specified language.
func (s *EmbeddedStore) Localize(hash string, lang string) string {
	if s.locs == nil {
		return hash
	}
	if langDict, ok := s.locs[lang]; ok {
		if text, found := langDict[hash]; found && text != "" {
			return text
		}
	}
	// Fallback to English dictionary if key is missing in target language
	if enDict, ok := s.locs["en"]; ok {
		if text, found := enDict[hash]; found && text != "" {
			return text
		}
	}
	return hash
}

func (s *EmbeddedStore) AvatarMeta(id int) (AvatarMeta, bool) {
	m, ok := s.avatars[id]
	return m, ok
}
func (s *EmbeddedStore) AvatarSkillsMeta(avatarID int) ([]SkillMeta, bool) {
	m, ok := s.skills[avatarID]
	return m, ok
}
func (s *EmbeddedStore) SkillTemplateMeta(skillID int) (SkillTemplateMeta, bool) {
	m, ok := s.skillTemplates[skillID]
	return m, ok
}
func (s *EmbeddedStore) AvatarMindscapesMeta(avatarID int) ([]MindscapeMeta, bool) {
	m, ok := s.mindscapes[avatarID]
	return m, ok
}
func (s *EmbeddedStore) AvatarPotentialVisionsMeta(avatarID int) ([]PotentialVisionMeta, bool) {
	m, ok := s.potentialVisions[avatarID]
	return m, ok
}
func (s *EmbeddedStore) WeaponMeta(id int) (WeaponMeta, bool) {
	m, ok := s.weapons[id]
	return m, ok
}
func (s *EmbeddedStore) EquipmentMeta(id int) (EquipmentMeta, bool) {
	m, ok := s.equipments[id]
	return m, ok
}
func (s *EmbeddedStore) EquipmentSuitMeta(suitID int) (EquipmentSuitMeta, bool) {
	m, ok := s.equipmentSuits[suitID]
	return m, ok
}
func (s *EmbeddedStore) PropertyMeta(id int) (PropertyMeta, bool) {
	m, ok := s.properties[id]
	return m, ok
}
func (s *EmbeddedStore) MedalMeta(id int) (MedalMeta, bool) {
	m, ok := s.medals[id]
	return m, ok
}
func (s *EmbeddedStore) TitleMeta(id int) (TitleMeta, bool) {
	m, ok := s.titles[id]
	return m, ok
}
func (s *EmbeddedStore) NamecardMeta(id int) (NamecardMeta, bool) {
	m, ok := s.namecards[id]
	return m, ok
}
func (s *EmbeddedStore) PfpMeta(id int) (PfpMeta, bool) {
	m, ok := s.pfps[id]
	return m, ok
}
func (s *EmbeddedStore) SkinMeta(id int) (SkinMeta, bool) {
	val, ok := s.skins[id]
	return val, ok
}

// DefaultSkinMeta returns the default skin metadata for an agent by their avatar ID.
func (s *EmbeddedStore) DefaultSkinMeta(avatarID int) (SkinMeta, bool) {
	skin, ok := s.defaultSkins[avatarID]
	return skin, ok
}

// WeaponLevelTemplate returns the multiplier template for a weapon's main stat at a given level.
func (s *EmbeddedStore) WeaponLevelTemplate(rarity, level int) (WeaponLevelTemplate, bool) {
	tpl, ok := s.weaponLevelTemplates[templateKey{Rarity: rarity, Level: level}]
	return tpl, ok
}

// WeaponStarTemplate returns the multiplier template for a weapon's stats based on refinement phase.
func (s *EmbeddedStore) WeaponStarTemplate(rarity, phase int) (WeaponStarTemplate, bool) {
	tpl, ok := s.weaponStarTemplates[templateKey{Rarity: rarity, Level: phase}]
	return tpl, ok
}

// EquipmentLevelTemplate returns the multiplier template for a drive disc's main stat.
func (s *EmbeddedStore) EquipmentLevelTemplate(rarity, level int) (EquipmentLevelTemplate, bool) {
	tpl, ok := s.equipLevelTemplates[templateKey{Rarity: rarity, Level: level}]
	return tpl, ok
}
