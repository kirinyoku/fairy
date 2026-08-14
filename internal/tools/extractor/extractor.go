package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func readJSON(path string, v interface{}) {
	b, err := os.ReadFile(path)
	check(err)
	err = json.Unmarshal(b, v)
	check(err)
}

func writeJSON(path string, v interface{}) {
	b, err := json.MarshalIndent(v, "", "  ")
	check(err)
	err = os.WriteFile(path, b, 0644)
	check(err)
}

func main() {
	assetsDir := filepath.Join("internal", "assets", "data")
	zenlessDir := "zenlessdata"

	// 1. Load targets
	var weapons map[string]map[string]interface{}
	readJSON(filepath.Join(assetsDir, "weapons.json"), &weapons)

	var equipments map[string]interface{}
	readJSON(filepath.Join(assetsDir, "equipments.json"), &equipments)

	var locs map[string]map[string]string
	readJSON(filepath.Join(assetsDir, "locs.json"), &locs)

	// 2. Load TextMaps for all supported languages
	langFiles := map[string]string{
		"de":    "TextMap_DETemplateTb.json",
		"en":    "TextMap_ENTemplateTb.json",
		"es":    "TextMap_ESTemplateTb.json",
		"fr":    "TextMap_FRTemplateTb.json",
		"id":    "TextMap_IDTemplateTb.json",
		"ja":    "TextMap_JATemplateTb.json",
		"ko":    "TextMap_KOTemplateTb.json",
		"pt":    "TextMap_PTTemplateTb.json",
		"ru":    "TextMap_RUTemplateTb.json",
		"th":    "TextMap_THTemplateTb.json",
		"vi":    "TextMap_VITemplateTb.json",
		"zh-cn": "TextMapTemplateTb.json",
		"zh-tw": "TextMap_CHTTemplateTb.json",
	}

	textMaps := make(map[string]map[string]string)
	for lang, filename := range langFiles {
		var tm map[string]string
		readJSON(filepath.Join(zenlessDir, "TextMap", filename), &tm)
		textMaps[lang] = tm
	}

	keysToExtract := make(map[string]bool)

	// 3. Process Weapons
	type WeaponTalentItem struct {
		WeaponID int    `json:"COEEBFOBGND"`
		Phase    int    `json:"APAEMLCPFID"`
		DescKey  string `json:"POLEJGCKKFI"`
	}
	var weaponTalentData struct {
		Items []WeaponTalentItem `json:"MLOEFHJHCID"`
	}
	readJSON(filepath.Join(zenlessDir, "FileCfg", "WeaponTalentTemplateTb.json"), &weaponTalentData)

	wtMap := make(map[string][]string)
	for _, item := range weaponTalentData.Items {
		wid := fmt.Sprintf("%d", item.WeaponID)
		if item.DescKey != "" {
			wtMap[wid] = append(wtMap[wid], item.DescKey)
		}
	}

	for id, w := range weapons {
		if keys, ok := wtMap[id]; ok && len(keys) > 0 {
			w["PassiveDescKeys"] = keys
			for _, k := range keys {
				keysToExtract[k] = true
			}
		} else {
			var keys []string
			for i := 1; i <= 5; i++ {
				key := fmt.Sprintf("Weapon_TalentDes_8%s0%d", id, i)
				keys = append(keys, key)
				keysToExtract[key] = true
			}
			w["PassiveDescKeys"] = keys
		}
	}
	writeJSON(filepath.Join(assetsDir, "weapons.json"), weapons)

	// 4. Process Equipments
	suitsRaw := equipments["Suits"].(map[string]interface{})
	for id, s := range suitsRaw {
		suit := s.(map[string]interface{})
		k2 := fmt.Sprintf("EquipmentSuit_%s_2_des", id)
		k4 := fmt.Sprintf("EquipmentSuit_%s_4_des", id)
		suit["Set2DescKey"] = k2
		suit["Set4DescKey"] = k4
		keysToExtract[k2] = true
		keysToExtract[k4] = true
	}
	writeJSON(filepath.Join(assetsDir, "equipments.json"), equipments)

	// 5. Process Skills
	type SkillDesItem struct {
		AvatarID int    `json:"PJABHBNCJOI"`
		Title    string `json:"LECKPHICFOA"`
		Desc     string `json:"DLADMENPFPD"`
		Formula  string `json:"KLPLBBJABBL"`
	}
	var skillDesData struct {
		Items []SkillDesItem `json:"MLOEFHJHCID"`
	}
	readJSON(filepath.Join(zenlessDir, "FileCfg", "AvatarSkillDesTemplateTb.json"), &skillDesData)

	type PassiveSkillDesItem struct {
		AvatarID int      `json:"PJABHBNCJOI"`
		Titles   []string `json:"OHBPFPBILBL"`
		Descs    []string `json:"IHDLHHGEHMP"`
	}
	var passiveSkillDesData struct {
		Items []PassiveSkillDesItem `json:"MLOEFHJHCID"`
	}
	readJSON(filepath.Join(zenlessDir, "FileCfg", "AvatarPassiveSkillDesTemplateTb.json"), &passiveSkillDesData)

	type SkillParamMeta struct {
		NameKey string `json:"NameKey"`
		Formula string `json:"Formula"`
	}

	type SkillMeta struct {
		NameKey string           `json:"NameKey"`
		DescKey string           `json:"DescKey"`
		Params  []SkillParamMeta `json:"Params,omitempty"`
	}

	skillsOutput := make(map[string][]SkillMeta)
	seenKeys := make(map[string]map[string]bool)
	currentSkillIdxMap := make(map[string]int)

	// Group skills and params per avatar
	for _, item := range skillDesData.Items {
		avStr := fmt.Sprintf("%d", item.AvatarID)
		if seenKeys[avStr] == nil {
			seenKeys[avStr] = make(map[string]bool)
		}

		if item.Title != "" {
			if seenKeys[avStr][item.Title] {
				continue
			}
			seenKeys[avStr][item.Title] = true

			skillsOutput[avStr] = append(skillsOutput[avStr], SkillMeta{
				NameKey: item.Title,
				DescKey: item.Desc,
				Params:  make([]SkillParamMeta, 0),
			})
			currentSkillIdxMap[avStr] = len(skillsOutput[avStr]) - 1
			keysToExtract[item.Title] = true
			keysToExtract[item.Desc] = true
		} else if item.Desc != "" && len(skillsOutput[avStr]) > 0 {
			// Check if desc matches a sub-skill title key (switching current skill section)
			matchingIdx := -1
			for idx, sk := range skillsOutput[avStr] {
				if sk.NameKey == item.Desc {
					matchingIdx = idx
					break
				}
			}
			if matchingIdx != -1 {
				currentSkillIdxMap[avStr] = matchingIdx
			} else if item.Formula != "" {
				cIdx := currentSkillIdxMap[avStr]
				skillsOutput[avStr][cIdx].Params = append(skillsOutput[avStr][cIdx].Params, SkillParamMeta{
					NameKey: item.Desc,
					Formula: item.Formula,
				})
				keysToExtract[item.Desc] = true
			}
		}
	}

	for _, item := range passiveSkillDesData.Items {
		avStr := fmt.Sprintf("%d", item.AvatarID)
		if seenKeys[avStr] == nil {
			seenKeys[avStr] = make(map[string]bool)
		}
		for i := 0; i < len(item.Titles); i++ {
			if item.Titles[i] == "" || seenKeys[avStr][item.Titles[i]] {
				continue
			}
			seenKeys[avStr][item.Titles[i]] = true

			skillsOutput[avStr] = append(skillsOutput[avStr], SkillMeta{
				NameKey: item.Titles[i],
				DescKey: item.Descs[i],
			})
			keysToExtract[item.Titles[i]] = true
			keysToExtract[item.Descs[i]] = true
		}
	}
	writeJSON(filepath.Join(assetsDir, "skills.json"), skillsOutput)

	// Extract skill templates (multiplier calculations)
	type SkillTemplateRaw struct {
		ID       int `json:"DALBKGGEJEF"`
		BaseDmg  int `json:"IKAABAIDFAO"`
		GrowDmg  int `json:"DGHHKAHHIPM"`
		BaseStun int `json:"OMFJHOLBIKA"`
		GrowStun int `json:"KICLLNBEAEN"`
		ExtraDmg int `json:"ECHPKCNANMI"`
		Cost     int `json:"BLGOMFMHNKA"`
	}

	type SkillTemplateMeta struct {
		BaseDmg  int `json:"bd"`
		GrowDmg  int `json:"gd"`
		BaseStun int `json:"bs"`
		GrowStun int `json:"gs"`
		ExtraDmg int `json:"ex"`
		Cost     int `json:"ec"`
	}

	var skillTemplateData struct {
		Items []SkillTemplateRaw `json:"MLOEFHJHCID"`
	}
	readJSON(filepath.Join(zenlessDir, "FileCfg", "AvatarSkillTemplateTb.json"), &skillTemplateData)

	skillTemplatesOutput := make(map[string]SkillTemplateMeta)
	for _, item := range skillTemplateData.Items {
		if item.ID == 0 {
			continue
		}
		idStr := fmt.Sprintf("%d", item.ID)
		skillTemplatesOutput[idStr] = SkillTemplateMeta{
			BaseDmg:  item.BaseDmg,
			GrowDmg:  item.GrowDmg,
			BaseStun: item.BaseStun,
			GrowStun: item.GrowStun,
			ExtraDmg: item.ExtraDmg,
			Cost:     item.Cost,
		}
	}
	writeJSON(filepath.Join(assetsDir, "skill_templates.json"), skillTemplatesOutput)

	// 5.5 Process Professions
	type ProfessionItem struct {
		NameKey string `json:"PCPJBDFHKNK"`
		DescKey string `json:"OLBIEFEPDLC"`
	}
	var professionData struct {
		Items []ProfessionItem `json:"MLOEFHJHCID"`
	}
	readJSON(filepath.Join(zenlessDir, "FileCfg", "AvatarProfessionTemplateTb.json"), &professionData)

	for _, item := range professionData.Items {
		if item.NameKey != "" {
			keysToExtract[item.NameKey] = true
		}
		if item.DescKey != "" {
			keysToExtract[item.DescKey] = true
		}
	}

	// 5.6 Process Weapon and Equipment Templates
	type rawWeaponLevelItem struct {
		Rarity             int `json:"ICPMKHFGPOG"`
		Level              int `json:"EMLFBEMHINK"`
		MainStat           int `json:"AHMDJCIHNKG"`
		SubStatDenominator int `json:"IDBKOAPHGLC"`
	}
	var rawWeaponLevelData struct {
		List []rawWeaponLevelItem `json:"OOFFGGKCDID"`
	}
	weaponLevelFile := filepath.Join(zenlessDir, "FileCfg", "WeaponLevelTemplateTb.json")
	if _, err := os.Stat(weaponLevelFile); err == nil {
		readJSON(weaponLevelFile, &rawWeaponLevelData)
		type cleanWeaponLevel struct {
			Rarity             int `json:"rarity"`
			Level              int `json:"level"`
			MainStat           int `json:"main_stat"`
			SubStatDenominator int `json:"sub_stat_denominator"`
		}
		var cleanData struct {
			List []cleanWeaponLevel `json:"list"`
		}
		for _, it := range rawWeaponLevelData.List {
			cleanData.List = append(cleanData.List, cleanWeaponLevel(it))
		}
		writeJSON(filepath.Join(assetsDir, "WeaponLevelTemplateTb.json"), cleanData)
	}

	type rawWeaponStarItem struct {
		Rarity     int `json:"ICPMKHFGPOG"`
		BreakLevel int `json:"BBOCBHBGMML"`
		MainStat   int `json:"NMFHJKEFLOG"`
		SubStat    int `json:"FCLIIPBDDKP"`
	}
	var rawWeaponStarData struct {
		List []rawWeaponStarItem `json:"OOFFGGKCDID"`
	}
	weaponStarFile := filepath.Join(zenlessDir, "FileCfg", "WeaponStarTemplateTb.json")
	if _, err := os.Stat(weaponStarFile); err == nil {
		readJSON(weaponStarFile, &rawWeaponStarData)
		type cleanWeaponStar struct {
			Rarity     int `json:"rarity"`
			BreakLevel int `json:"break_level"`
			MainStat   int `json:"main_stat"`
			SubStat    int `json:"sub_stat"`
		}
		var cleanData struct {
			List []cleanWeaponStar `json:"list"`
		}
		for _, it := range rawWeaponStarData.List {
			cleanData.List = append(cleanData.List, cleanWeaponStar(it))
		}
		writeJSON(filepath.Join(assetsDir, "WeaponStarTemplateTb.json"), cleanData)
	}

	type rawEquipLevelItem struct {
		Rarity   int `json:"GMKDLJLLBPO"`
		Level    int `json:"FNPIELBFDEJ"`
		MainStat int `json:"JEKGLLBALFE"`
	}
	var rawEquipLevelData struct {
		List []rawEquipLevelItem `json:"MIJCMCEDADM"`
	}
	equipLevelFile := filepath.Join(zenlessDir, "FileCfg", "EquipmentLevelTemplateTb.json")
	if _, err := os.Stat(equipLevelFile); err == nil {
		readJSON(equipLevelFile, &rawEquipLevelData)
		type cleanEquipLevel struct {
			Rarity   int `json:"rarity"`
			Level    int `json:"level"`
			MainStat int `json:"main_stat"`
		}
		var cleanData struct {
			List []cleanEquipLevel `json:"list"`
		}
		for _, it := range rawEquipLevelData.List {
			cleanData.List = append(cleanData.List, cleanEquipLevel(it))
		}
		writeJSON(filepath.Join(assetsDir, "EquipmentLevelTemplateTb.json"), cleanData)
	}

	// 6. Update Locs
	addedStats := make(map[string]int)
	for lang := range langFiles {
		if locs[lang] == nil {
			locs[lang] = make(map[string]string)
		}
	}

	for k := range keysToExtract {
		for lang, tm := range textMaps {
			if val, ok := tm[k]; ok {
				locs[lang][k] = val
				addedStats[lang]++
			}
		}
	}

	writeJSON(filepath.Join(assetsDir, "locs.json"), locs)
	fmt.Printf("Extraction complete. Processed %d keys across %d languages.\n", len(keysToExtract), len(langFiles))
	for lang, count := range addedStats {
		fmt.Printf("  - %s: %d strings\n", lang, count)
	}
}
