package sibylValues

import (
	"errors"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

var (
	HelperBot         *gotgbot.Bot
	BotUpdater        *ext.Updater
	SendReportHandler ReportHandler
)

var (
	ErrInvalidPerm = errors.New("invalid permission provided")
)

/*
Coefficients and Flags

==== Flags     -  ========
Range 0-100 (No bans) (Dominator Locked)
• Civilian     - 0-80
• Past Banned  - 81-100
==============
Range 100-300 (Auto-mute) (Non-lethal Paralyzer)
• TROLLING     - 101-125 - Trolling
• SPAM         - 126-200 - Spam/Unwanted Promotion
• EVADE        - 201-250 - Ban Evade using alts
x-------x
Manual Revert
• CUSTOM       - 251-300 - Any Custom reason
x-------x
==============
Range 300+ (Ban on Sight) (Lethal Eliminator)
• PSYCHOHAZARD - 301-350 - Bulk banned due to some bad users
• MALIMP       - 351-400 - Malicious Impersonation
• NSFW         - 401-450 - Sending NSFW Content in SFW
• RAID         - 451-500 - Bulk join raid to vandalize
• MASSADD      - 501-600 - Mass adding to group/channel
==============
*/

// crime coefficient increasement ranges
var (
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
)
