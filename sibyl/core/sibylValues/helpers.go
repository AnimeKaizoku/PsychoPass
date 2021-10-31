package sibylValues

import (
	"strconv"
	"strings"

	"github.com/ALiwoto/StrongStringGo/strongStringGo"
	"github.com/AnimeKaizoku/PsychoPass/sibyl/core/utils/timeUtils"
)

func NewReport(reason, message string, target, reporter int64,
	reporterPerm UserPermission) *Report {

	return &Report{
		ReportReason:       reason,
		ReportMessage:      message,
		TargetUser:         target,
		ReporterId:         reporter,
		ReportDate:         timeUtils.GenerateCurrentDateTime(),
		ReporterPermission: reporterPerm.GetStringPermission(),
	}
}

func ConvertToPermission(value string) (UserPermission, error) {
	value = strings.ToLower(strings.TrimSpace(value))
	// first of all check and see if value is an integer or not
	valueInt, err := strconv.Atoi(value)
	if err == nil {
		perm := UserPermission(valueInt)
		if perm.IsValid() {
			return perm, nil
		}
		// we already know that the value is a valid integer, so there is no
		// chance that the value is a valid permission in string format.
		return NormalUser, ErrInvalidPerm
	}

	switch value {
	case "user", "civilian":
		return NormalUser, nil
	case "enforcer":
		return Enforcer, nil
	case "inspector":
		return Inspector, nil
	case "owner":
		return Owner, nil
	default:
		return NormalUser, ErrInvalidPerm
	}
}

func IsInvalidID(id int64) bool {
	return id == 777000 || id == 1087968824
}

func GetCrimeCoefficientRange(value int) *CrimeCoefficientRange {
	/*
		RangeCivilian     = &CrimeCoefficientRange{0, 80}
		RangePastBanned   = &CrimeCoefficientRange{81, 100}
		RangeTrolling     = &CrimeCoefficientRange{101, 125}
		RangeSpam         = &CrimeCoefficientRange{126, 200}
		RangeEvade        = &CrimeCoefficientRange{201, 250}
		RangeCustom       = &CrimeCoefficientRange{251, 300}
		RangePsychoHazard = &CrimeCoefficientRange{301, 350}
		RangeMalImp       = &CrimeCoefficientRange{351, 400}
		RangeNSFW         = &CrimeCoefficientRange{401, 450}
		RangeRaid         = &CrimeCoefficientRange{451, 500}
		RangeMassAdd      = &CrimeCoefficientRange{501, 600}
	*/
	if value < 0 {
		return nil
	} else if RangeCivilian.IsInRange(value) {
		return RangeCivilian
	} else if RangePastBanned.IsInRange(value) {
		return RangePastBanned
	} else if RangeTrolling.IsInRange(value) {
		return RangeTrolling
	} else if RangeSpam.IsInRange(value) {
		return RangeSpam
	} else if RangeEvade.IsInRange(value) {
		return RangeEvade
	} else if RangeCustom.IsInRange(value) {
		return RangeCustom
	} else if RangePsychoHazard.IsInRange(value) {
		return RangePsychoHazard
	} else if RangeMalImp.IsInRange(value) {
		return RangeMalImp
	} else if RangeNSFW.IsInRange(value) {
		return RangeNSFW
	} else if RangeRaid.IsInRange(value) {
		return RangeRaid
	} else if RangeMassAdd.IsInRange(value) {
		return RangeMassAdd
	}

	return nil
}

func GetCCRangeByString(value string) []*CrimeCoefficientRange {
	value = fixReasonString(strings.ToLower(strings.TrimSpace(value)))
	values := strongStringGo.Split(value, " ", "\n", ",", "|", "\t", ";",
		".", "-", "--")
	var tmp *CrimeCoefficientRange
	var result []*CrimeCoefficientRange
	exists := func(c *CrimeCoefficientRange) bool {
		if len(result) == 0 {
			return false
		}

		for _, v := range result {
			if v != nil && v.IsValueInRange(c) {
				return true
			}
		}
		return false
	}

	for _, current := range values {
		tmp = getCCRangeByString(current)
		if tmp != nil && !exists(tmp) {
			result = append(result, tmp)
		}
	}

	if len(result) == 0 {
		result = append(result, RangeCustom)
	}
	return result
}

func getCCRangeByString(value string) *CrimeCoefficientRange {
	/*
		// Range 0-100 (No bans) (Dominator Locked)
		// Civilian     - 0-80
		// Past Banned  - 81-100
		// Range 100-300 (Auto-mute) (Non-lethal Paralyzer)
		ReasonTrolling = "trolling"
		ReasonSpam     = "spam"
		ReasonEvade    = "evade"
		ReasonCustom   = "evade"
		// Range 300+ (Ban on Sight) (Lethal Eliminator)
		ReasonMalimp       = "malimp"
		ReasonPsychoHazard = "psychohazard"
		ReasonNSFW         = "nsfw"
		ReasonRaid         = "raid"
		ReasonMassAdd      = "massadd"
	*/
	if strings.HasPrefix(value, ReasonTrolling) {
		return RangeTrolling
	} else if strings.HasPrefix(value, ReasonSpam) {
		return RangeSpam
	} else if strings.HasPrefix(value, ReasonEvade) {
		return RangeEvade
	} else if strings.HasPrefix(value, ReasonMalImp) {
		return RangeMalImp
	} else if strings.HasPrefix(value, ReasonPsychoHazard) {
		return RangePsychoHazard
	} else if strings.HasPrefix(value, ReasonNSFW) {
		return RangeNSFW
	} else if strings.HasPrefix(value, ReasonRaid) {
		return RangeRaid
	} else if strings.HasPrefix(value, ReasonMassAdd) {
		return RangeMassAdd
	}

	return nil
}

func fixReasonString(value string) string {
	value = strings.ReplaceAll(value, "mass add", ReasonMassAdd)
	value = strings.ReplaceAll(value, "mass-add", ReasonMassAdd)
	value = strings.ReplaceAll(value, "n.s.f.w", ReasonNSFW)
	value = strings.ReplaceAll(value, "psycho hazard", ReasonPsychoHazard)
	value = strings.ReplaceAll(value, "psycho-hazard", ReasonPsychoHazard)
	return value
}
