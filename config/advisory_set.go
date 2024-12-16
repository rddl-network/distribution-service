package config

import "github.com/rddl-network/distribution-service/model"

func GetWeeklyAdvisoryDistribution() [1]model.Advisory {
	return [1]model.Advisory{
		{
			Address: "VJLHinV6iAVSw7Mwx1yRe3jLT86pbjoygJRPNroXhcKxAmtX2EZZM4wAhW993umuquWG7wujcPXw98f9",
			Amount:  9615.3846154,
		},
	}
}
