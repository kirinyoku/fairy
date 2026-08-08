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
	assetsDir := filepath.Join("..", "..", "assets", "data")
	zenlessDir := filepath.Join("..", "..", "..", "zenlessdata")

	// 1. Load targets
	var weapons map[string]map[string]interface{}
	readJSON(filepath.Join(assetsDir, "weapons.json"), &weapons)

	var equipments map[string]interface{}
	readJSON(filepath.Join(assetsDir, "equipments.json"), &equipments)

	var locs map[string]map[string]string
	readJSON(filepath.Join(assetsDir, "locs.json"), &locs)

	// 2. Load TextMaps
	var enText, ruText map[string]string
	readJSON(filepath.Join(zenlessDir, "TextMap", "TextMap_ENTemplateTb.json"), &enText)
	readJSON(filepath.Join(zenlessDir, "TextMap", "TextMap_RUTemplateTb.json"), &ruText)

	keysToExtract := make(map[string]bool)

	// 3. Process Weapons
	for id, w := range weapons {
		var keys []string
		for i := 1; i <= 5; i++ {
			key := fmt.Sprintf("Weapon_TalentDes_8%s0%d", id, i)
			keys = append(keys, key)
			keysToExtract[key] = true
		}
		w["PassiveDescKeys"] = keys
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

	// 6. Update Locs
	if locs["en"] == nil {
		locs["en"] = make(map[string]string)
	}
	if locs["ru"] == nil {
		locs["ru"] = make(map[string]string)
	}

	addedEn := 0
	addedRu := 0
	for k := range keysToExtract {
		if val, ok := enText[k]; ok {
			locs["en"][k] = val
			addedEn++
		}
		if val, ok := ruText[k]; ok {
			locs["ru"][k] = val
			addedRu++
		}
	}

	writeJSON(filepath.Join(assetsDir, "locs.json"), locs)
	fmt.Printf("Extraction complete. Added %d EN strings and %d RU strings.\n", addedEn, addedRu)
}
