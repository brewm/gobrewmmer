package main

import (
	"log"
	"strings"

	"github.com/brewm/gobrewmmer/api/recepie"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

func main() {
	recepieStruct := &recepie.Recepie{
		Uuid:        1,
		Name:        "Little Cub IPA - Bronze Cub",
		Description: "Brown IPA",
		Ingredients: []*recepie.Ingredient{
			{
				Type: recepie.IngredientType_WATER,
				Name: "Init water",
				Quantity: &recepie.Quantity{
					Volume: 32,
					Unit:   recepie.Unit_L,
				},
			},
			{
				Type: recepie.IngredientType_WATER,
				Name: "Sparge water",
				Quantity: &recepie.Quantity{
					Volume: 25,
					Unit:   recepie.Unit_L,
				},
			},
			{
				Type: recepie.IngredientType_MALT,
				Name: "Pale Ale",
				Quantity: &recepie.Quantity{
					Volume: 8,
					Unit:   recepie.Unit_KG,
				},
			},
			{
				Type: recepie.IngredientType_MALT,
				Name: "Crystal",
				Quantity: &recepie.Quantity{
					Volume: 1.5,
					Unit:   recepie.Unit_KG,
				},
			},
			{
				Type: recepie.IngredientType_HOPS,
				Name: "Tomahawk",
				Quantity: &recepie.Quantity{
					Volume: 75,
					Unit:   recepie.Unit_G,
				},
			},
			{
				Type: recepie.IngredientType_HOPS,
				Name: "Cascade",
				Quantity: &recepie.Quantity{
					Volume: 80,
					Unit:   recepie.Unit_G,
				},
			},
			{
				Type: recepie.IngredientType_HOPS,
				Name: "Citra",
				Quantity: &recepie.Quantity{
					Volume: 170,
					Unit:   recepie.Unit_G,
				},
			},
			{
				Type: recepie.IngredientType_YIEST,
				Name: "Safale US-05",
				Quantity: &recepie.Quantity{
					Volume: 20,
					Unit:   recepie.Unit_G,
				},
			}},
		Steps: []*recepie.Step{
			{
				Phase:       recepie.Phase_INIT,
				Temperature: 55,
				Ingredients: []*recepie.Ingredient{
					{
						Type: recepie.IngredientType_WATER,
						Name: "Init water",
						Quantity: &recepie.Quantity{
							Volume: 32,
							Unit:   recepie.Unit_L,
						},
					},
					{
						Type: recepie.IngredientType_MALT,
						Name: "Pale Ale",
						Quantity: &recepie.Quantity{
							Volume: 8,
							Unit:   recepie.Unit_KG,
						},
					},
					{
						Type: recepie.IngredientType_MALT,
						Name: "Crystal",
						Quantity: &recepie.Quantity{
							Volume: 1.5,
							Unit:   recepie.Unit_KG,
						},
					},
				},
			},
			{
				Phase:       recepie.Phase_MASHING,
				Temperature: 55,
				Duration: &recepie.Quantity{
					Volume: 30,
					Unit:   recepie.Unit_MIN,
				},
			},
			{
				Phase:       recepie.Phase_MASHING,
				Temperature: 65,
				Duration: &recepie.Quantity{
					Volume: 45,
					Unit:   recepie.Unit_MIN,
				},
			},
			{
				Phase:       recepie.Phase_MASHING,
				Temperature: 78,
				Duration: &recepie.Quantity{
					Volume: 45,
					Unit:   recepie.Unit_MIN,
				},
			},
			{
				Phase:       recepie.Phase_SPARGING,
				Temperature: 55,
				Ingredients: []*recepie.Ingredient{
					{
						Type: recepie.IngredientType_WATER,
						Name: "Sparge water",
						Quantity: &recepie.Quantity{
							Volume: 25,
							Unit:   recepie.Unit_L,
						},
					},
				},
			},
			{
				Phase: recepie.Phase_BOILING,
				Ingredients: []*recepie.Ingredient{
					{
						Type: recepie.IngredientType_HOPS,
						Name: "Tomahawk",
						Quantity: &recepie.Quantity{
							Volume: 75,
							Unit:   recepie.Unit_G,
						},
					},
				},
				Duration: &recepie.Quantity{
					Volume: 60,
					Unit:   recepie.Unit_MIN,
				},
			},
			{
				Phase: recepie.Phase_BOILING,
				Ingredients: []*recepie.Ingredient{
					{
						Type: recepie.IngredientType_HOPS,
						Name: "Cascade",
						Quantity: &recepie.Quantity{
							Volume: 80,
							Unit:   recepie.Unit_G,
						},
					},
				},
				Duration: &recepie.Quantity{
					Volume: 10,
					Unit:   recepie.Unit_MIN,
				},
			},
			{
				Phase: recepie.Phase_BOILING,
				Ingredients: []*recepie.Ingredient{
					{
						Type: recepie.IngredientType_HOPS,
						Name: "Citra",
						Quantity: &recepie.Quantity{
							Volume: 75,
							Unit:   recepie.Unit_G,
						},
					},
				},
			},
			{
				Phase: recepie.Phase_FERMENTATION,
				Duration: &recepie.Quantity{
					Volume: 7,
					Unit:   recepie.Unit_DAY,
				},
			},
			{
				Phase: recepie.Phase_FERMENTATION,
				Ingredients: []*recepie.Ingredient{
					{
						Type: recepie.IngredientType_HOPS,
						Name: "Citra",
						Quantity: &recepie.Quantity{
							Volume: 95,
							Unit:   recepie.Unit_G,
						},
					},
				},
				Duration: &recepie.Quantity{
					Volume: 7,
					Unit:   recepie.Unit_DAY,
				},
			},
		},
	}

	data, err := proto.Marshal(recepieStruct)
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}
	serialized := &recepie.Recepie{}
	err = proto.Unmarshal(data, serialized)
	if err != nil {
		log.Fatal("unmarshaling error: ", err)
	}

	m := jsonpb.Marshaler{}
	recepieJson, err := m.MarshalToString(recepieStruct)
	if err != nil {
		log.Fatal("json marshaling error: ", err)
	}
	println(recepieJson)

	um := jsonpb.Unmarshaler{}
	unserialized := &recepie.Recepie{}
	err = um.Unmarshal(strings.NewReader(recepieJson), unserialized)
	if err != nil {
		log.Fatal("json unmarshaling error: ", err)
	}
	println(unserialized)

	recepieJson2, err := m.MarshalToString(unserialized)
	if err != nil {
		log.Fatal("json marshaling error: ", err)
	}
	println(recepieJson2)

}
