package fairy

import (
	"math"
	"strings"

	"github.com/kirinyoku/fairy/store"
)

// calcAgentBaseStat calculates the base stats of an agent from level growth and core enhancements.
// ZZZ Math Quirk:
// Agent Base stats consist of four parts:
// 1. Initial base prop at level 1 (BaseProps from AvatarBaseTemplateTb)
// 2. Growth from leveling up, divided by 10000 (GrowthProps from AvatarBaseTemplateTb)
// 3. Flat additions from Promotion / Ascensions (PromotionProps from AvatarPromotionTemplateTb)
// 4. Flat additions from Core Skill Enhancements (CoreEnhancementProps from AvatarSkillCoreTemplateTb)
//
// The result should always be rounded down (math.Floor) to match the game's display.
func calcAgentBaseStat(meta store.AvatarMeta, propID, level, promotionLevel, coreEnhancement int) float64 {
	baseVal, _ := meta.BaseStat(propID)
	base := float64(baseVal)

	// GrowthValue = (GrowthProps[PropertyId] * (Avatar.Level - 1)) / 10000
	growthVal, _ := meta.GrowthStat(propID)
	growth := float64(growthVal*(level-1)) / 10000.0

	val := base + growth

	// 3. Promotion prop
	promVal, _ := meta.PromotionStat(promotionLevel, propID)
	val += float64(promVal)

	// 4. Core enhancement prop
	coreVal, _ := meta.CoreEnhancementStat(coreEnhancement, propID)
	val += float64(coreVal)

	return val
}

// calcWEngineMainStat calculates the main stat value of a weapon.
// W-Engine Main Stat uses two multipliers that are ADDITIVE to each other:
// 1. The level multiplier from WeaponLevelTemplate (AHMDJCIHNKG)
// 2. The phase/refinement multiplier from WeaponStarTemplate (NMFHJKEFLOG)
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
	// Result = MainStat.PropertyValue * (1 + WeaponLevel.FIELD_XXX / 10000 + WeaponStar.FIELD_YYY / 10000)
	result := float64(baseVal) * (1.0 + float64(levelMod)/10000.0 + float64(starMod)/10000.0)
	return int(math.Floor(result))
}

// calcWEngineSecondaryStat calculates the actual secondary stat value of a weapon.
// W-Engine Secondary Stats are completely undocumented in API wrappers, but they DO scale with level.
// However, instead of a direct level multiplier, they use a "Denominator" from WeaponLevelTemplate (IDBKOAPHGLC).
// This is an inverse scaling trick. The level multiplier is calculated as `10000 / IDBKOAPHGLC`.
// Example: At level 60, IDBKOAPHGLC is 4000, which gives a 2.5x multiplier to the base secondary stat (10000 / 4000 = 2.5).
func calcWEngineSecondaryStat(s store.MetadataStore, meta store.WeaponMeta, level, phase int) int {
	baseVal := meta.SecondaryStat.PropertyValue
	levelMult := 1.0
	if lvlTpl, ok := s.WeaponLevelTemplate(meta.Rarity, level); ok {
		if lvlTpl.SubStatDenominator > 0 {
			levelMult = 10000.0 / float64(lvlTpl.SubStatDenominator)
		}
	}

	starMod := 0
	if starTpl, ok := s.WeaponStarTemplate(meta.Rarity, phase); ok {
		starMod = starTpl.SubStat
	}

	result := float64(baseVal) * levelMult * (1.0 + float64(starMod)/10000.0)
	// Removed debug print
	return int(math.Floor(result))
}

// calculateAgentStats populates the BaseStats and Stats of an Agent.
// It executes the full stat pipeline:
// 1. Calculate base stats from character level, promotion, and core skills.
// 2. Add W-Engine base attack (which merges into the character's base attack).
// 3. Accumulate all percentage multipliers and flat bonuses from W-Engine substats, drive discs, and set bonuses.
// 4. Apply the bonuses to the base stats in a strict order (Base * (1 + PercentBonus) + FlatBonus) to compute the final stats.
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
	baseEnergyRegen := calcAgentBaseStat(meta, int(PropBaseEnergyRegen), agent.Level, agent.Promotion, agent.CoreSkillEnhancement) / 100.0
	baseSheerForce := calcAgentBaseStat(meta, int(PropBaseSheerForce), agent.Level, agent.Promotion, agent.CoreSkillEnhancement)

	// Fixed Base stats
	baseCritRate := 0.05
	baseCritDMG := 0.50
	basePenRatio := 0.0
	basePenFlat := 0.0

	// W-Engine Base ATK is accumulated separately because it is added to the Agent's
	// Base ATK *before* percentage multipliers are applied.
	wEngineBaseAtk := 0.0

	// For stats calculation, we aggregate percent multipliers and flat bonuses.
	// We accumulate all sources of bonuses (W-Engine substats, Drive Discs, Set Bonuses)
	// into this map first, and then apply them all at once at the end.
	bonuses := make(map[int]float64)

	addBonus := func(propID int, value float64) {
		bonuses[propID] += value
	}

	// Calculate W-Engine stats
	if agent.WEngine != nil {
		if wMeta, ok := s.WeaponMeta(agent.WEngine.ID); ok {
			// Modification is 1-indexed (1-5), so phase is Modification - 1 (0-4)
			wPhase := agent.WEngine.Modification - 1
			if wPhase < 0 {
				wPhase = 0
			}

			wMainStatId := wMeta.MainStat.PropertyID
			wMainStatVal := calcWEngineMainStat(s, wMeta, agent.WEngine.Level, agent.WEngine.Phase)

			// Main stat of weapon is usually ATK (12101)
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
				addBonus(wSecStatId, float64(wSecStatVal)/10000.0)
			} else {
				addBonus(wSecStatId, float64(wSecStatVal))
			}
		}
	}

	baseAtk = math.Floor(baseAtk) + math.Floor(wEngineBaseAtk)
	baseDef = math.Floor(baseDef)
	baseImpact = math.Floor(baseImpact)
	baseAnomalyMastery = math.Floor(baseAnomalyMastery)
	baseAnomalyProficiency = math.Floor(baseAnomalyProficiency)
	basePenFlat = math.Floor(basePenFlat)
	baseSheerForce = math.Floor(baseSheerForce)
	if meta.ProfessionType == "Rupture" || agent.Specialty == SpecialtyRupture {
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

	// Add Drive Disc stats
	// 1. Disc Main Stat MUST be multiplied by the Level Modifier from EquipmentLevelTemplateTb.json.
	//    The raw JSON value is just the Level 0 base stat.
	// 2. Disc Sub Stats are provided as the base value per roll.
	//    They MUST be multiplied by the number of `Rolls` the stat received during upgrades.
	for _, disc := range agent.DriveDiscs {
		if disc.MainStat.PropertyID != 0 {
			addBonus(int(disc.MainStat.PropertyID), disc.MainStat.Value)
		}
		for _, sub := range disc.SubStats {
			addBonus(int(sub.PropertyID), sub.Value)
		}
	}

	// Add Set Bonuses
	for _, bonus := range agent.ActiveSetBonuses {
		if suitMeta, ok := s.EquipmentSuitMeta(bonus.Set.ID); ok {
			for propID, val := range suitMeta.SetBonusProps {
				isPercent := false
				if pMeta, pOk := s.PropertyMeta(propID); pOk {
					if strings.Contains(pMeta.Format, "%") {
						isPercent = true
					}
				}
				if isPercent {
					addBonus(propID, float64(val)/10000.0)
				} else {
					addBonus(propID, float64(val))
				}
			}
		}
	}

	// Apply bonuses in order of operations.
	// 1. Multiply the aggregated Base Stat by (1 + sum of all Percent Multipliers).
	// 2. Add the sum of all Flat Bonuses.
	// 3. Most stats are floored, but CritRate, CritDMG, PenRatio, and EnergyRegen are NOT floored.
	totalHp := math.Floor(baseHp*(1.0+bonuses[int(PropHPPercent)]+bonuses[int(PropHPPercentBonus)]) + bonuses[int(PropHPFlat)] + bonuses[int(PropHPFlatBonus)])
	totalAtk := math.Floor(baseAtk*(1.0+bonuses[int(PropATKPercent)]) + bonuses[int(PropATKFlat)])

	totalSheerForce := math.Floor(bonuses[int(PropBaseSheerForce)] + bonuses[int(PropSheerForce)])
	if meta.ProfessionType == "Rupture" || agent.Specialty == SpecialtyRupture {
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
		EnergyRegen:        baseEnergyRegen*(1.0+bonuses[int(PropEnergyRegenPercent)]) + bonuses[int(PropBaseEnergyRegen)] + bonuses[int(PropEnergyRegen)],
		SheerForce:         totalSheerForce,
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
