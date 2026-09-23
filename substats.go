package fairy

// agentRecommendedSubStats maps specific agent IDs to their verified in-game recommended Drive Disc sub-stats.
//
// Data Source & In-Game Verification:
// These recommendations are neither arbitrary nor subjective theorycrafting; they are extracted directly
// from the official Zenless Zone Zero in-game equipment recommendation system (the yellow highlight indicators
// shown on Drive Disc sub-stats in the game's equipment, tuning, and "Recommend" / L3 build guide screens).
var agentRecommendedSubStats = map[int][]PropertyID{
	// 1011: Anby (Stun) - ATK%, CRIT Rate, CRIT DMG
	1011: {PropATKPercent, PropCritRate, PropCritDMG},

	// 1021: Nekomata (Attack) - ATK%, CRIT Rate, CRIT DMG
	1021: {PropATKPercent, PropCritRate, PropCritDMG},

	// 1031: Nicole (Support) - ATK%, Anomaly Proficiency
	1031: {PropATKPercent, PropAnomalyProficiency},

	// 1041: Soldier 11 (Attack) - ATK%, CRIT Rate, CRIT DMG
	1041: {PropATKPercent, PropCritRate, PropCritDMG},

	// 1051: Yidhari (Rupture) - HP%, CRIT Rate, CRIT DMG
	1051: {PropHPPercent, PropCritRate, PropCritDMG},

	// 1061: Corin (Attack) - ATK%, CRIT Rate, CRIT DMG
	1061: {PropATKPercent, PropCritRate, PropCritDMG},

	// 1071: Caesar (Defense) - ATK%, CRIT Rate, CRIT DMG
	1071: {PropATKPercent, PropCritRate, PropCritDMG},

	// 1081: Billy (Attack) - ATK%, CRIT Rate, CRIT DMG
	1081: {PropATKPercent, PropCritRate, PropCritDMG},

	// 1091: Miyabi (Anomaly) - ATK%, CRIT Rate, CRIT DMG
	1091: {PropATKPercent, PropCritRate, PropCritDMG},

	// 1101: Koleda (Stun) - ATK%, CRIT Rate, CRIT DMG
	1101: {PropATKPercent, PropCritRate, PropCritDMG},

	// 1111: Anton (Attack) - ATK%, CRIT Rate, CRIT DMG
	1111: {PropATKPercent, PropCritRate, PropCritDMG},

	// 1121: Ben (Defense) - DEF%, ATK%
	1121: {PropDEFPercent, PropATKPercent},

	// 1131: Soukaku (Support) - ATK%, Anomaly Proficiency
	1131: {PropATKPercent, PropAnomalyProficiency},

	// 1141: Lycaon (Stun) - ATK%, CRIT Rate, CRIT DMG
	1141: {PropATKPercent, PropCritRate, PropCritDMG},

	// 1151: Lucy (Support) - ATK%, Flat ATK
	1151: {PropATKPercent, PropATKFlat},

	// 1161: Lighter (Stun) - ATK%, CRIT Rate, CRIT DMG
	1161: {PropATKPercent, PropCritRate, PropCritDMG},

	// 1171: Burnice (Anomaly) - ATK%, Anomaly Proficiency
	1171: {PropATKPercent, PropAnomalyProficiency},

	// 1181: Grace (Anomaly) - ATK%, Anomaly Proficiency
	1181: {PropATKPercent, PropAnomalyProficiency},

	// 1191: Ellen (Attack) - ATK%, CRIT Rate, CRIT DMG
	1191: {PropATKPercent, PropCritRate, PropCritDMG},

	// 1201: Harumasa (Attack) - ATK%, CRIT Rate, CRIT DMG
	1201: {PropATKPercent, PropCritRate, PropCritDMG},

	// 1211: Rina (Support) - ATK%, Anomaly Proficiency
	1211: {PropATKPercent, PropAnomalyProficiency},

	// 1221: Yanagi (Anomaly) - ATK%, Anomaly Proficiency
	1221: {PropATKPercent, PropAnomalyProficiency},

	// 1241: Zhu Yuan (Attack) - ATK%, CRIT Rate, CRIT DMG
	1241: {PropATKPercent, PropCritRate, PropCritDMG},

	// 1251: Qingyi (Stun) - ATK%, CRIT Rate, CRIT DMG
	1251: {PropATKPercent, PropCritRate, PropCritDMG},

	// 1261: Jane Doe (Anomaly) - ATK%, Anomaly Proficiency
	1261: {PropATKPercent, PropAnomalyProficiency},

	// 1271: Seth (Defense) - ATK%, Anomaly Proficiency
	1271: {PropATKPercent, PropAnomalyProficiency},

	// 1281: Piper (Anomaly) - ATK%, Anomaly Proficiency
	1281: {PropATKPercent, PropAnomalyProficiency},

	// 1291: Hugo Vlad (Attack) - ATK%, CRIT Rate, CRIT DMG
	1291: {PropATKPercent, PropCritRate, PropCritDMG},

	// 1301: Orphie & Magus (Attack) - ATK%, CRIT Rate, CRIT DMG
	1301: {PropATKPercent, PropCritRate, PropCritDMG},

	// 1311: Astra Yao (Support) - ATK%, Flat ATK
	1311: {PropATKPercent, PropATKFlat},

	// 1321: Evelyn (Attack) - ATK%, CRIT Rate, CRIT DMG
	1321: {PropATKPercent, PropCritRate, PropCritDMG},

	// 1331: Vivian (Anomaly) - ATK%, Anomaly Proficiency
	1331: {PropATKPercent, PropAnomalyProficiency},

	// 1341: Zhao (Defense) - HP%, Flat HP
	1341: {PropHPPercent, PropHPFlat},

	// 1351: Pulchra (Stun) - ATK%, CRIT Rate, CRIT DMG
	1351: {PropATKPercent, PropCritRate, PropCritDMG},

	// 1361: Trigger (Attack) - ATK%, CRIT Rate, CRIT DMG
	1361: {PropATKPercent, PropCritRate, PropCritDMG},

	// 1371: Yixuan (Rupture) - HP%, CRIT Rate, CRIT DMG
	1371: {PropHPPercent, PropCritRate, PropCritDMG},

	// 1381: Soldier 0 - Anby (Attack) - ATK%, CRIT Rate, CRIT DMG
	1381: {PropATKPercent, PropCritRate, PropCritDMG},

	// 1391: Ju Fufu (Attack) - ATK%, CRIT Rate, CRIT DMG
	1391: {PropATKPercent, PropCritRate, PropCritDMG},

	// 1401: Alice (Anomaly) - ATK%, Anomaly Proficiency
	1401: {PropATKPercent, PropAnomalyProficiency},

	// 1411: Yuzuha (Anomaly) - ATK%, Anomaly Proficiency
	1411: {PropATKPercent, PropAnomalyProficiency},

	// 1421: Pan Yinhu (Defense) - ATK%, Flat ATK
	1421: {PropATKPercent, PropATKFlat},

	// 1431: Ye Shunguang (Attack) - ATK%, CRIT Rate, CRIT DMG
	1431: {PropATKPercent, PropCritRate, PropCritDMG},

	// 1441: Manato (Rupture) - HP%, CRIT Rate, CRIT DMG
	1441: {PropHPPercent, PropCritRate, PropCritDMG},

	// 1451: Lucia (Support) - HP%, Flat HP
	1451: {PropHPPercent, PropHPFlat},

	// 1461: Seed (Attack) - ATK%, CRIT Rate, CRIT DMG
	1461: {PropATKPercent, PropCritRate, PropCritDMG},

	// 1471: Banyue (Rupture) - HP%, CRIT Rate, CRIT DMG
	1471: {PropHPPercent, PropCritRate, PropCritDMG},

	// 1481: Dialyn (Attack) - ATK%, CRIT Rate, CRIT DMG
	1481: {PropATKPercent, PropCritRate, PropCritDMG},

	// 1491: Sunna (Support) - ATK%, Flat ATK
	1491: {PropATKPercent, PropATKFlat},

	// 1501: Aria (Anomaly) - ATK%, Anomaly Proficiency
	1501: {PropATKPercent, PropAnomalyProficiency},

	// 1511: Nangong Yu (Stun) - ATK%, Anomaly Proficiency
	1511: {PropATKPercent, PropAnomalyProficiency},

	// 1521: Cissia (Attack) - ATK%, CRIT Rate, CRIT DMG
	1521: {PropATKPercent, PropCritRate, PropCritDMG},

	// 1531: Starlight - Billy (Rupture) - HP%, CRIT Rate, CRIT DMG
	1531: {PropHPPercent, PropCritRate, PropCritDMG},

	// 1541: Promeia (Anomaly) - ATK%, Anomaly Proficiency
	1541: {PropATKPercent, PropAnomalyProficiency},

	// 1551: Pyrois (Attack) - ATK%, CRIT Rate, CRIT DMG
	1551: {PropATKPercent, PropCritRate, PropCritDMG},

	// 1561: Velina (Anomaly) - ATK%, Anomaly Proficiency
	1561: {PropATKPercent, PropAnomalyProficiency},

	// 1571: Norma (Stun) - ATK%, CRIT Rate, CRIT DMG
	1571: {PropATKPercent, PropCritRate, PropCritDMG},

	// 1581: Remielle (Anomaly) - ATK%, Flat ATK, Anomaly Proficiency
	1581: {PropATKPercent, PropATKFlat, PropAnomalyProficiency},

	// 1591: Sigrid (Attack) - ATK%, CRIT Rate, CRIT DMG
	1591: {PropATKPercent, PropCritRate, PropCritDMG},

	// 1611: Claret (Armorer) - DEF%, CRIT Rate, CRIT DMG
	1611: {PropDEFPercent, PropCritRate, PropCritDMG},

	// 1621: Roxy (Stun) - ATK%, CRIT Rate, CRIT DMG
	1621: {PropATKPercent, PropCritRate, PropCritDMG},
}

// AgentRecommendedSubStats returns the recommended Drive Disc sub-stat property IDs for the given agent ID,
// using the shared default client.
//
// These stats reflect the official in-game recommendation system in Zenless Zone Zero, matching the yellow
// highlight indicators displayed on Drive Disc sub-stats in the equipment and tuning UI for each specific agent.
//
// If a curated list exists in [agentRecommendedSubStats], a defensive copy is returned.
// Otherwise, it computes a safe fallback from [AgentHighlightProps], converting attributes to percentage-only
// sub-stats and excluding flat ATK, flat HP, flat DEF, and flat PEN.
func AgentRecommendedSubStats(agentID int) []PropertyID {
	client, err := getDefaultClient()
	if err != nil {
		return []PropertyID{}
	}
	return client.AgentRecommendedSubStats(agentID)
}

// AllAgentRecommendedSubStats returns a map of all known agents to their recommended Drive Disc sub-stats,
// matching the yellow sub-stat highlights shown in the in-game equipment UI, using the shared default client.
func AllAgentRecommendedSubStats() map[int][]PropertyID {
	client, err := getDefaultClient()
	if err != nil {
		return nil
	}
	return client.AllAgentRecommendedSubStats()
}
