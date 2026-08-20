package fairy

import (
	"math"
	"strings"

	"github.com/kirinyoku/fairy/internal/store"
)

// calcAgentBaseStat calculates the base stats of an Agent from level growth, Promotion phase, and Core Skill Enhancement.
//
// In Zenless Zone Zero, an Agent's base stats before gear consist of four components:
//  1. Initial base stat at level 1 (BaseProps in AvatarBaseTemplateTb).
//  2. Growth from leveling up, scaled by statModifierScale (10000):
//     GrowthValue = (GrowthProps[PropertyId] * (Agent.Level - 1)) / 10000.
//  3. Flat stat additions from Promotion / Ascension phase (PromotionProps in AvatarPromotionTemplateTb).
//  4. Flat stat additions from Core Skill Enhancements (CoreEnhancementProps in AvatarSkillCoreTemplateTb).
func calcAgentBaseStat(meta store.AvatarMeta, propID, level, promotionLevel, coreEnhancement int) float64 {
	baseVal, _ := meta.BaseStat(propID)
	base := float64(baseVal)

	// GrowthValue = (GrowthProps[PropertyId] * (Avatar.Level - 1)) / 10000
	growthVal, _ := meta.GrowthStat(propID)
	growth := float64(growthVal*(level-1)) / statModifierScale

	val := base + growth

	// 3. Promotion prop
	promVal, _ := meta.PromotionStat(promotionLevel, propID)
	val += float64(promVal)

	// 4. Core enhancement prop
	coreVal, _ := meta.CoreEnhancementStat(coreEnhancement, propID)
	val += float64(coreVal)

	return val
}

// calcWEngineMainStat calculates the main stat value (Base ATK) of a W-Engine.
//
// W-Engine Main Stat uses two additive multipliers:
//  1. The level multiplier from WeaponLevelTemplate (AHMDJCIHNKG / 10000).
//  2. The star/phase multiplier from WeaponStarTemplate (NMFHJKEFLOG / 10000).
//
// Formula: BaseValue * (1 + LevelMod/10000 + StarMod/10000) rounded down.
func calcWEngineMainStat(s store.MetadataStore, meta store.WeaponMeta, level, phase int) int {
	baseVal := meta.MainStat.PropertyValue
	levelMod := 0
	if lvlTpl, ok := s.WeaponLevelTemplate(meta.Rarity, level); ok {
		levelMod = lvlTpl.MainStat
	}
	starMod := 0
	if starTpl, ok := s.WeaponStarTemplate(meta.Rarity, phase); ok {
		starMod = starTpl.MainStat
	}
	// Result = MainStat.PropertyValue * (1 + WeaponLevel.MainStat / 10000 + WeaponStar.MainStat / 10000)
	result := float64(baseVal) * (1.0 + float64(levelMod)/statModifierScale + float64(starMod)/statModifierScale)
	return int(math.Floor(result))
}

// calcWEngineSecondaryStat calculates the level-scaled secondary stat value of a W-Engine.
//
// W-Engine secondary stats scale via an inverse denominator from WeaponLevelTemplate (SubStatDenominator):
//
//	LevelMultiplier = 10000 / SubStatDenominator
//
// For example, at level 60 with denominator 4000, LevelMultiplier is 2.5x the base secondary stat (10000 / 4000 = 2.5).
func calcWEngineSecondaryStat(s store.MetadataStore, meta store.WeaponMeta, level, phase int) int {
	baseVal := meta.SecondaryStat.PropertyValue
	levelMult := 1.0
	if lvlTpl, ok := s.WeaponLevelTemplate(meta.Rarity, level); ok {
		if lvlTpl.SubStatDenominator > 0 {
			levelMult = statModifierScale / float64(lvlTpl.SubStatDenominator)
		}
	}

	starMod := 0
	if starTpl, ok := s.WeaponStarTemplate(meta.Rarity, phase); ok {
		starMod = starTpl.SubStat
	}

	result := float64(baseVal) * levelMult * (1.0 + float64(starMod)/statModifierScale)
	return int(math.Floor(result))
}

// calculateAgentStats populates the BaseStats and final Stats of an [Agent].
//
// Pipeline execution order:
//  1. Calculate innate base stats from Agent level growth, Promotion phase, and Core Skill Enhancements.
//  2. Add W-Engine Base ATK (which merges directly into the Agent's innate Base ATK).
//  3. Accumulate all percentage multipliers and flat bonuses from W-Engine substats, Drive Discs, and set bonuses.
//  4. Compute final combat stats applying multipliers: Base * (1 + PercentBonus) + FlatBonus.
//  5. Apply ZZZ rounding rules: math.Floor for HP, ATK, DEF, Impact, Anomaly, PenFlat, SheerForce;
//     exact floating-point decimals for CritRate, CritDMG, PenRatio, and EnergyRegen.
func calculateAgentStats(agent *Agent, s store.MetadataStore) {
	meta, ok := s.AvatarMeta(agent.ID)
	if !ok {
		return
	}

	// 1. Calculate Base Stats (Avatar + W-Engine Base ATK)
	baseHp := calcAgentBaseStat(meta, int(PropBaseHP), agent.Level, agent.Promotion, agent.CoreSkillEnhancement)
	baseAtk := calcAgentBaseStat(meta, int(PropBaseATK), agent.Level, agent.Promotion, agent.CoreSkillEnhancement)
	baseDef := calcAgentBaseStat(meta, int(PropBaseDEF), agent.Level, agent.Promotion, agent.CoreSkillEnhancement)
	baseImpact := calcAgentBaseStat(meta, int(PropBaseImpact), agent.Level, agent.Promotion, agent.CoreSkillEnhancement)
	baseAnomalyMastery := calcAgentBaseStat(meta, int(PropBaseAnomalyMastery), agent.Level, agent.Promotion, agent.CoreSkillEnhancement)
	baseAnomalyProficiency := calcAgentBaseStat(meta, int(PropBaseAnomalyProficiency), agent.Level, agent.Promotion, agent.CoreSkillEnhancement)
	baseEnergyRegenProp := PropBaseEnergyRegen
	if agent.Specialty == SpecialtyRupture {
		baseEnergyRegenProp = PropBaseRpRecover
	}
	baseEnergyRegen := calcAgentBaseStat(meta, int(baseEnergyRegenProp), agent.Level, agent.Promotion, agent.CoreSkillEnhancement) / 100.0
	baseSheerForce := calcAgentBaseStat(meta, int(PropBaseSheerForce), agent.Level, agent.Promotion, agent.CoreSkillEnhancement)

	// Fixed Base stats
	baseCritRate := defaultBaseCritRate
	baseCritDMG := defaultBaseCritDMG
	basePenRatio := defaultBasePenRatio
	basePenFlat := defaultBasePenFlat

	bonuses := make(map[int]float64)
	addBonus := func(propID int, value float64) {
		bonuses[propID] += value
	}

	wEngineBaseAtk := accumulateWEngineBonus(agent, s, addBonus)

	baseAtk = math.Floor(baseAtk) + math.Floor(wEngineBaseAtk)
	baseDef = math.Floor(baseDef)
	baseImpact = math.Floor(baseImpact)
	baseAnomalyMastery = math.Floor(baseAnomalyMastery)
	baseAnomalyProficiency = math.Floor(baseAnomalyProficiency)
	basePenFlat = math.Floor(basePenFlat)
	baseSheerForce = math.Floor(baseSheerForce)
	if agent.Specialty == SpecialtyRupture {
		baseSheerForce += math.Floor(baseHp*0.1) + math.Floor(baseAtk*0.3)
	}

	agent.BaseStats = Stats{
		HP:                 math.Floor(baseHp),
		ATK:                baseAtk,
		DEF:                baseDef,
		Impact:             baseImpact,
		CritRate:           baseCritRate,
		CritDMG:            baseCritDMG,
		AnomalyMastery:     baseAnomalyMastery,
		AnomalyProficiency: baseAnomalyProficiency,
		PenRatio:           basePenRatio,
		PenFlat:            basePenFlat,
		EnergyRegen:        baseEnergyRegen,
		SheerForce:         baseSheerForce,
	}

	accumulateDriveDiscBonus(agent, addBonus)
	accumulateSetBonus(agent, s, addBonus)

	// Apply bonuses in order of operations.
	// 1. Multiply the aggregated Base Stat by (1 + sum of all Percent Multipliers).
	// 2. Add the sum of all Flat Bonuses.
	// 3. Most stats are floored, but CritRate, CritDMG, PenRatio, and EnergyRegen are NOT floored.
	totalHp := math.Floor(baseHp*(1.0+bonuses[int(PropHPPercent)]+bonuses[int(PropHPPercentBonus)]) + bonuses[int(PropHPFlat)] + bonuses[int(PropHPFlatBonus)])
	totalAtk := math.Floor(baseAtk*(1.0+bonuses[int(PropATKPercent)]) + bonuses[int(PropATKFlat)])

	totalSheerForce := math.Floor(bonuses[int(PropBaseSheerForce)] + bonuses[int(PropSheerForce)])
	if agent.Specialty == SpecialtyRupture {
		totalSheerForce += math.Floor(totalHp*0.1) + math.Floor(totalAtk*0.3)
	}

	// Calculate Attribute DMG Bonus matching the agent's elemental attribute
	totalAttrDMG := sumPropVariants(bonuses, propGroupGeneralDMG)
	switch agent.Attribute.BaseAttribute() {
	case AttributePhysical:
		totalAttrDMG += sumPropVariants(bonuses, propGroupPhysicalDMG)
	case AttributeFire:
		totalAttrDMG += sumPropVariants(bonuses, propGroupFireDMG)
	case AttributeIce:
		totalAttrDMG += sumPropVariants(bonuses, propGroupIceDMG)
	case AttributeElectric:
		totalAttrDMG += sumPropVariants(bonuses, propGroupElectricDMG)
	case AttributeEther:
		totalAttrDMG += sumPropVariants(bonuses, propGroupEtherDMG)
	case AttributeWind:
		totalAttrDMG += sumPropVariants(bonuses, propGroupWindDMG)
	}

	agent.Stats = Stats{
		HP:                 totalHp,
		ATK:                totalAtk,
		DEF:                math.Floor(baseDef*(1.0+bonuses[int(PropDEFPercent)]) + bonuses[int(PropDEFFlat)]),
		Impact:             math.Floor(baseImpact*(1.0+bonuses[int(PropImpactPercent)]) + bonuses[int(PropImpactFlat)]),
		CritRate:           baseCritRate + bonuses[int(PropBaseCritRate)] + bonuses[int(PropCritRate)],
		CritDMG:            baseCritDMG + bonuses[int(PropBaseCritDMG)] + bonuses[int(PropCritDMG)],
		AttributeDMGBonus:  totalAttrDMG,
		AnomalyMastery:     math.Floor(baseAnomalyMastery*(1.0+bonuses[int(PropAnomalyMasteryPercent)]) + bonuses[int(PropBaseAnomalyMastery)] + bonuses[int(PropAnomalyMastery)]),
		AnomalyProficiency: math.Floor(baseAnomalyProficiency*(1.0+bonuses[int(PropAnomalyProficiencyPercent)]) + bonuses[int(PropBaseAnomalyProficiency)] + bonuses[int(PropAnomalyProficiency)]),
		PenRatio:           basePenRatio + bonuses[int(PropBasePENRatio)] + bonuses[int(PropPENRatio)],
		PenFlat:            math.Floor(basePenFlat + bonuses[int(PropBasePENFlat)] + bonuses[int(PropPENFlat)]),
		EnergyRegen:        baseEnergyRegen*(1.0+bonuses[int(PropEnergyRegenPercent)]+bonuses[int(PropRpRecoverPercent)]) + bonuses[int(PropBaseEnergyRegen)] + bonuses[int(PropEnergyRegen)] + bonuses[int(PropBaseRpRecover)] + bonuses[int(PropRpRecover)],
		SheerForce:         totalSheerForce,
	}
}

func accumulateWEngineBonus(agent *Agent, s store.MetadataStore, addBonus func(int, float64)) float64 {
	var wEngineBaseAtk float64
	if agent.WEngine == nil {
		return wEngineBaseAtk
	}
	wMeta, ok := s.WeaponMeta(agent.WEngine.ID)
	if !ok {
		return wEngineBaseAtk
	}

	wPhase := agent.WEngine.Modification - 1
	if wPhase < 0 {
		wPhase = 0
	}

	wMainStatId := wMeta.MainStat.PropertyID
	wMainStatVal := calcWEngineMainStat(s, wMeta, agent.WEngine.Level, agent.WEngine.Phase)

	if wMainStatId == int(PropBaseATK) {
		wEngineBaseAtk += float64(wMainStatVal)
	}

	wSecStatId := wMeta.SecondaryStat.PropertyID
	wSecStatVal := calcWEngineSecondaryStat(s, wMeta, agent.WEngine.Level, wPhase)

	isPercent := false
	if pMeta, pOk := s.PropertyMeta(wSecStatId); pOk {
		if strings.Contains(pMeta.Format, "%") {
			isPercent = true
		}
	}

	if isPercent {
		addBonus(wSecStatId, float64(wSecStatVal)/statModifierScale)
	} else {
		addBonus(wSecStatId, float64(wSecStatVal))
	}
	return wEngineBaseAtk
}

func accumulateDriveDiscBonus(agent *Agent, addBonus func(int, float64)) {
	for _, disc := range agent.DriveDiscs.Slots {
		if disc.MainStat.PropertyID != 0 {
			addBonus(int(disc.MainStat.PropertyID), disc.MainStat.Value)
		}
		for _, sub := range disc.SubStats {
			addBonus(int(sub.PropertyID), sub.Value)
		}
	}
}

func accumulateSetBonus(agent *Agent, s store.MetadataStore, addBonus func(int, float64)) {
	for _, bonus := range agent.DriveDiscs.SetBonuses {
		if bonus.Count < 2 {
			continue
		}
		if suitMeta, ok := s.EquipmentSuitMeta(int(bonus.Set.ID)); ok {
			for propID, val := range suitMeta.SetBonusProps {
				isPercent := false
				if pMeta, pOk := s.PropertyMeta(propID); pOk {
					if strings.Contains(pMeta.Format, "%") {
						isPercent = true
					}
				}
				if isPercent {
					addBonus(propID, float64(val)/statModifierScale)
				} else {
					addBonus(propID, float64(val))
				}
			}
		}
	}
}

// sumPropVariants sums all property variations (base, percent, flat, percent bonus, flat bonus)
// for a given base property group ID.
func sumPropVariants(bonuses map[int]float64, baseGroup PropertyID) float64 {
	var total float64
	baseID := int(baseGroup)
	for i := 1; i <= 5; i++ {
		total += bonuses[baseID+i]
	}
	return total
}
