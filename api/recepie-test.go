package main

import (
	"log"
	"strings"

	"github.com/brewm/gobrewmmer/pkg/api/brewmmer"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

func main() {
	recepieStruct := &brewmmer.Recepie{
		Id:          1,
		Name:        "Little Cub IPA - Bronze Cub",
		Description: "Brown IPA",
		Ingredients: []*brewmmer.Ingredient{
			{
				Type: brewmmer.IngredientType_WATER,
				Name: "Init water",
				Quantity: &brewmmer.Quantity{
					Volume: 32,
					Unit:   brewmmer.Unit_L,
				},
			},
			{
				Type: brewmmer.IngredientType_WATER,
				Name: "Sparge water",
				Quantity: &brewmmer.Quantity{
					Volume: 25,
					Unit:   brewmmer.Unit_L,
				},
			},
			{
				Type: brewmmer.IngredientType_MALT,
				Name: "Pale Ale",
				Quantity: &brewmmer.Quantity{
					Volume: 8,
					Unit:   brewmmer.Unit_KG,
				},
			},
			{
				Type: brewmmer.IngredientType_MALT,
				Name: "Crystal",
				Quantity: &brewmmer.Quantity{
					Volume: 1.5,
					Unit:   brewmmer.Unit_KG,
				},
			},
			{
				Type: brewmmer.IngredientType_HOPS,
				Name: "Tomahawk",
				Quantity: &brewmmer.Quantity{
					Volume: 75,
					Unit:   brewmmer.Unit_G,
				},
			},
			{
				Type: brewmmer.IngredientType_HOPS,
				Name: "Cascade",
				Quantity: &brewmmer.Quantity{
					Volume: 80,
					Unit:   brewmmer.Unit_G,
				},
			},
			{
				Type: brewmmer.IngredientType_HOPS,
				Name: "Citra",
				Quantity: &brewmmer.Quantity{
					Volume: 170,
					Unit:   brewmmer.Unit_G,
				},
			},
			{
				Type: brewmmer.IngredientType_YIEST,
				Name: "Safale US-05",
				Quantity: &brewmmer.Quantity{
					Volume: 20,
					Unit:   brewmmer.Unit_G,
				},
			}},
		Steps: []*brewmmer.Step{
			{
				Phase:       brewmmer.Phase_INIT,
				Temperature: 55,
				Ingredients: []*brewmmer.Ingredient{
					{
						Type: brewmmer.IngredientType_WATER,
						Name: "Init water",
						Quantity: &brewmmer.Quantity{
							Volume: 32,
							Unit:   brewmmer.Unit_L,
						},
					},
					{
						Type: brewmmer.IngredientType_MALT,
						Name: "Pale Ale",
						Quantity: &brewmmer.Quantity{
							Volume: 8,
							Unit:   brewmmer.Unit_KG,
						},
					},
					{
						Type: brewmmer.IngredientType_MALT,
						Name: "Crystal",
						Quantity: &brewmmer.Quantity{
							Volume: 1.5,
							Unit:   brewmmer.Unit_KG,
						},
					},
				},
			},
			{
				Phase:       brewmmer.Phase_MASHING,
				Temperature: 55,
				Duration: &brewmmer.Quantity{
					Volume: 30,
					Unit:   brewmmer.Unit_MIN,
				},
			},
			{
				Phase:       brewmmer.Phase_MASHING,
				Temperature: 65,
				Duration: &brewmmer.Quantity{
					Volume: 45,
					Unit:   brewmmer.Unit_MIN,
				},
			},
			{
				Phase:       brewmmer.Phase_MASHING,
				Temperature: 78,
				Duration: &brewmmer.Quantity{
					Volume: 45,
					Unit:   brewmmer.Unit_MIN,
				},
			},
			{
				Phase:       brewmmer.Phase_SPARGING,
				Temperature: 55,
				Ingredients: []*brewmmer.Ingredient{
					{
						Type: brewmmer.IngredientType_WATER,
						Name: "Sparge water",
						Quantity: &brewmmer.Quantity{
							Volume: 25,
							Unit:   brewmmer.Unit_L,
						},
					},
				},
			},
			{
				Phase: brewmmer.Phase_BOILING,
				Ingredients: []*brewmmer.Ingredient{
					{
						Type: brewmmer.IngredientType_HOPS,
						Name: "Tomahawk",
						Quantity: &brewmmer.Quantity{
							Volume: 75,
							Unit:   brewmmer.Unit_G,
						},
					},
				},
				Duration: &brewmmer.Quantity{
					Volume: 60,
					Unit:   brewmmer.Unit_MIN,
				},
			},
			{
				Phase: brewmmer.Phase_BOILING,
				Ingredients: []*brewmmer.Ingredient{
					{
						Type: brewmmer.IngredientType_HOPS,
						Name: "Cascade",
						Quantity: &brewmmer.Quantity{
							Volume: 80,
							Unit:   brewmmer.Unit_G,
						},
					},
				},
				Duration: &brewmmer.Quantity{
					Volume: 10,
					Unit:   brewmmer.Unit_MIN,
				},
			},
			{
				Phase: brewmmer.Phase_BOILING,
				Ingredients: []*brewmmer.Ingredient{
					{
						Type: brewmmer.IngredientType_HOPS,
						Name: "Citra",
						Quantity: &brewmmer.Quantity{
							Volume: 75,
							Unit:   brewmmer.Unit_G,
						},
					},
				},
			},
			{
				Phase: brewmmer.Phase_FERMENTATION,
				Duration: &brewmmer.Quantity{
					Volume: 7,
					Unit:   brewmmer.Unit_DAY,
				},
			},
			{
				Phase: brewmmer.Phase_FERMENTATION,
				Ingredients: []*brewmmer.Ingredient{
					{
						Type: brewmmer.IngredientType_HOPS,
						Name: "Citra",
						Quantity: &brewmmer.Quantity{
							Volume: 95,
							Unit:   brewmmer.Unit_G,
						},
					},
				},
				Duration: &brewmmer.Quantity{
					Volume: 7,
					Unit:   brewmmer.Unit_DAY,
				},
			},
		},
	}

	data, err := proto.Marshal(recepieStruct)
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}
	serialized := &brewmmer.Recepie{}
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
	unserialized := &brewmmer.Recepie{}
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
