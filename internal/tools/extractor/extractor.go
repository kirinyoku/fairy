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

	type SkillMeta struct {
		NameKey string `json:"NameKey"`
		DescKey string `json:"DescKey"`
	}

	skillsOutput := make(map[string][]SkillMeta)

	for _, item := range skillDesData.Items {
		if item.Title == "" {
			continue
		}
		avStr := fmt.Sprintf("%d", item.AvatarID)
		skillsOutput[avStr] = append(skillsOutput[avStr], SkillMeta{
			NameKey: item.Title,
			DescKey: item.Desc,
		})
		keysToExtract[item.Title] = true
		keysToExtract[item.Desc] = true
	}

	for _, item := range passiveSkillDesData.Items {
		avStr := fmt.Sprintf("%d", item.AvatarID)
		for i := 0; i < len(item.Titles); i++ {
			if item.Titles[i] == "" {
				continue
			}
			skillsOutput[avStr] = append(skillsOutput[avStr], SkillMeta{
				NameKey: item.Titles[i],
				DescKey: item.Descs[i],
			})
			keysToExtract[item.Titles[i]] = true
			keysToExtract[item.Descs[i]] = true
		}
	}
	writeJSON(filepath.Join(assetsDir, "skills.json"), skillsOutput)

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
