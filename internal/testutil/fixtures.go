package testutil

import "encoding/json"

// CompanyFixture returns a realistic JSON response for a company lookup.
func CompanyFixture() json.RawMessage {
	return json.RawMessage(`{
		"siren": "443061841",
		"denomination": "GOOGLE FRANCE",
		"nom_entreprise": "GOOGLE FRANCE",
		"forme_juridique": "SAS, société par actions simplifiée",
		"capital": 1000000,
		"date_creation": "2002-01-30",
		"siege": {
			"siret": "44306184100047",
			"adresse_ligne_1": "8 RUE DE LONDRES",
			"code_postal": "75009",
			"ville": "PARIS"
		},
		"dirigeants": [
			{
				"nom": "FITZPATRICK",
				"prenom": "Sebastien",
				"qualite": "Président"
			}
		]
	}`)
}

// SearchResultFixture returns a realistic search response with one result.
func SearchResultFixture() json.RawMessage {
	return json.RawMessage(`{
		"resultats": [
			{
				"siren": "443061841",
				"nom_entreprise": "GOOGLE FRANCE",
				"denomination": "GOOGLE FRANCE"
			}
		],
		"total": 1
	}`)
}

// EmptySearchFixture returns a search response with no results.
func EmptySearchFixture() json.RawMessage {
	return json.RawMessage(`{"resultats": [], "total": 0}`)
}

// AssociationFixture returns a realistic JSON response for an association.
func AssociationFixture() json.RawMessage {
	return json.RawMessage(`{
		"siren": "123456789",
		"numero_rna": "W751234567",
		"titre": "ASSOCIATION TEST",
		"objet": "Promotion de la technologie"
	}`)
}

// SuggestionsFixture returns a realistic suggestions response.
func SuggestionsFixture() json.RawMessage {
	return json.RawMessage(`[
		{
			"denomination": "GOOGLE FRANCE",
			"siren": "443061841",
			"forme_juridique": "SAS"
		}
	]`)
}

// CreditsFixture returns a realistic API credits response.
func CreditsFixture() json.RawMessage {
	return json.RawMessage(`{
		"jetons_utilises": 150,
		"jetons_restants": 850,
		"plan": "premium"
	}`)
}
