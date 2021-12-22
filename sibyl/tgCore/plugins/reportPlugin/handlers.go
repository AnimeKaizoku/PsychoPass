package reportPlugin

import (
	"strings"

	"github.com/ALiwoto/StrongStringGo/strongStringGo"
	"github.com/MinistryOfWelfare/PsychoPass/sibyl/core/sibylConfig"
	sv "github.com/MinistryOfWelfare/PsychoPass/sibyl/core/sibylValues"
	"github.com/MinistryOfWelfare/PsychoPass/sibyl/database"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func sendReportHandler(r *sv.Report) {
	// prevent from panic xD
	if sv.HelperBot == nil {
		return
	}

	bases := sibylConfig.GetBaseChatIds()
	if len(bases) == 0 {
		// there is no chat to send the report to...
		// ignore the report...
		return
	}

	var text string
	var opts *gotgbot.SendMessageOpts

	md := r.ParseAsMd()

	text = md.ToString()
	opts = &gotgbot.SendMessageOpts{
		ParseMode:             sv.MarkDownV2,
		DisableWebPagePreview: true,
		ReplyMarkup:           getReportButtons(r.UniqueId),
	}

	for _, chat := range bases {
		sendReportMessage(chat, text, opts)
	}
}

func sendMultiReportHandler(r *sv.MultiScanRawData) {
	// prevent from panic xD
	if sv.HelperBot == nil {
		return
	}

	bases := sibylConfig.GetBaseChatIds()
	if len(bases) == 0 {
		// there is no chat to send the report to...
		// ignore the report...
		return
	}

	var text string
	var opts *gotgbot.SendMessageOpts

	md := r.ParseAsMd()

	text = md.ToString()
	opts = &gotgbot.SendMessageOpts{
		ParseMode:             sv.MarkDownV2,
		DisableWebPagePreview: true,
		ReplyMarkup:           getMultiReportButtons(r.AssociationBanId),
	}

	for _, chat := range bases {
		sendReportMessage(chat, text, opts)
	}
}

func sendToADHandler(d *sv.AssaultDominatorData) {
	// prevent from panic xD
	if sv.HelperBot == nil {
		return
	}

	bases := sibylConfig.GetADIds()
	if len(bases) == 0 {
		// there is no chat to send the report to...
		// ignore the report...
		return
	}

	var opts *gotgbot.SendMessageOpts

	text := d.ParseAsText()
	opts = &gotgbot.SendMessageOpts{
		ParseMode:             sv.MarkDownV2,
		DisableWebPagePreview: true,
	}

	for _, chat := range bases {
		sendReportMessage(chat, text, opts)
	}
}

func scanCallBackQuery(cq *gotgbot.CallbackQuery) bool {
	return strings.HasPrefix(cq.Data, ReportPrefix+sepChar)
}

func scanCallBackResponse(b *gotgbot.Bot, ctx *ext.Context) error {
	tgUser := ctx.EffectiveUser
	token, err := database.GetTokenFromId(tgUser.Id)
	message := ctx.EffectiveMessage

	if err != nil || token == nil || !token.CanBan() {
		_, _ = ctx.CallbackQuery.Answer(b, &gotgbot.AnswerCallbackQueryOpts{
			Text:      "This is not for you...",
			ShowAlert: true,
			CacheTime: 5,
		})
		return ext.EndGroups
	}

	myStrs := strings.Split(ctx.CallbackQuery.Data, sepChar)
	// re_r_caseID
	// re_a_caseID
	// 0 _ 1 _ 2
	if len(myStrs) != 3 {
		return ext.ContinueGroups
	}

	data := myStrs[1]
	caseId := myStrs[2]
	scan := database.GetScan(caseId)
	if scan == nil {
		_, _ = ctx.CallbackQuery.Answer(b, &gotgbot.AnswerCallbackQueryOpts{
			Text:      "This scan is invalid...",
			ShowAlert: true,
			CacheTime: 5,
		})
		return ext.EndGroups
	}

	if !scan.IsPending() {
		_, _ = ctx.CallbackQuery.Answer(b, &gotgbot.AnswerCallbackQueryOpts{
			Text:      "This scan has already been " + scan.GetStatusString() + "...",
			ShowAlert: true,
			CacheTime: 5,
		})
		_, _ = message.EditText(b, scan.ParseAsMd().ToString(), &gotgbot.EditMessageTextOpts{
			DisableWebPagePreview: true,
		})
		return ext.EndGroups
	}

	scan.AgentUser = tgUser

	if data == CloseData {
		scan.Close(token.UserId, "") /* no reason */
		database.UpdateScan(scan)
		_, _ = message.EditText(b, scan.ParseAsMd().ToString(), &gotgbot.EditMessageTextOpts{
			ParseMode:             sv.MarkDownV2,
			DisableWebPagePreview: true,
		})
		_, _ = message.Delete(b)
		return ext.EndGroups
	}

	switch data {
	case ApproveData:
		scan.Approve(token.UserId, "") /* no reason */
		go pushScanToDatabase(scan)
	case RejectData:
		scan.Reject(token.UserId, "") /* no reason */
	}

	database.UpdateScan(scan)

	_, _ = message.EditText(b, scan.ParseAsMd().ToString(), &gotgbot.EditMessageTextOpts{
		ParseMode:             sv.MarkDownV2,
		DisableWebPagePreview: true,
	})

	return ext.EndGroups
}

func approveHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	message := ctx.EffectiveMessage
	replied := message.ReplyToMessage
	if message.ReplyToMessage == nil || replied.ReplyMarkup == nil {
		return ext.EndGroups
	}

	markup := replied.ReplyMarkup

	if len(markup.InlineKeyboard) == 0 || len(markup.InlineKeyboard[0]) == 0 {
		return ext.EndGroups
	}

	button := markup.InlineKeyboard[0][0]
	if button.CallbackData == "" {
		return ext.EndGroups
	}

	tgUser := ctx.EffectiveUser
	token, err := database.GetTokenFromId(tgUser.Id)
	if err != nil || token == nil || !token.CanBan() {
		return ext.EndGroups
	}

	query := button.CallbackData
	myStrs := strings.Split(query, sepChar)
	// re_r_caseID
	// re_a_caseID
	// 0 _ 1 _ 2
	if len(myStrs) != 3 {
		return ext.EndGroups
	}

	caseId := myStrs[2]
	scan := database.GetScan(caseId)
	if scan == nil {
		return ext.EndGroups
	}

	if !scan.IsPending() {
		_, _ = replied.EditText(b, scan.ParseAsMd().ToString(), &gotgbot.EditMessageTextOpts{
			DisableWebPagePreview: true,
		})
		return ext.EndGroups
	}

	args := strongStringGo.SplitN(message.Text, 2, " ", "\n")
	var newReason string
	if len(args) > 1 {
		newReason = args[1]
	}

	scan.AgentUser = tgUser

	scan.Approve(token.UserId, newReason)
	go pushScanToDatabase(scan)
	database.UpdateScan(scan)
	_, _ = replied.EditText(b, scan.ParseAsMd().ToString(), &gotgbot.EditMessageTextOpts{
		ParseMode:             sv.MarkDownV2,
		DisableWebPagePreview: true,
	})

	return ext.EndGroups
}

func rejectHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	message := ctx.EffectiveMessage
	replied := message.ReplyToMessage
	if message.ReplyToMessage == nil || replied.ReplyMarkup == nil {
		return ext.EndGroups
	}

	markup := replied.ReplyMarkup

	if len(markup.InlineKeyboard) == 0 || len(markup.InlineKeyboard[0]) == 0 {
		return ext.EndGroups
	}

	button := markup.InlineKeyboard[0][0]
	if button.CallbackData == "" {
		return ext.EndGroups
	}

	tgUser := ctx.EffectiveUser
	token, err := database.GetTokenFromId(tgUser.Id)
	if err != nil || token == nil || !token.CanBan() {
		return ext.EndGroups
	}

	query := button.CallbackData
	myStrs := strings.Split(query, sepChar)
	// re_r_caseID
	// re_a_caseID
	// 0 _ 1 _ 2
	if len(myStrs) != 3 {
		return ext.EndGroups
	}

	caseId := myStrs[2]
	scan := database.GetScan(caseId)
	if scan == nil {
		return ext.EndGroups
	}

	if !scan.IsPending() {
		_, _ = replied.EditText(b, scan.ParseAsMd().ToString(), &gotgbot.EditMessageTextOpts{
			DisableWebPagePreview: true,
		})
		return ext.EndGroups
	}

	args := strongStringGo.SplitN(message.Text, 2, " ", "\n")
	var newReason string
	if len(args) > 1 {
		newReason = args[1]
	}

	scan.AgentUser = tgUser

	scan.Reject(token.UserId, newReason)
	database.UpdateScan(scan)
	_, _ = replied.EditText(b, scan.ParseAsMd().ToString(), &gotgbot.EditMessageTextOpts{
		ParseMode:             sv.MarkDownV2,
		DisableWebPagePreview: true,
	})

	return ext.EndGroups
}

func closeHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	message := ctx.EffectiveMessage
	replied := message.ReplyToMessage
	if message.ReplyToMessage == nil || replied.ReplyMarkup == nil {
		return ext.EndGroups
	}

	markup := replied.ReplyMarkup

	if len(markup.InlineKeyboard) == 0 || len(markup.InlineKeyboard[0]) == 0 {
		return ext.EndGroups
	}

	button := markup.InlineKeyboard[0][0]
	if button.CallbackData == "" {
		return ext.EndGroups
	}

	tgUser := ctx.EffectiveUser
	token, err := database.GetTokenFromId(tgUser.Id)
	if err != nil || token == nil || !token.CanBan() {
		return ext.EndGroups
	}

	query := button.CallbackData
	myStrs := strings.Split(query, sepChar)
	// re_r_caseID
	// re_a_caseID
	// 0 _ 1 _ 2
	if len(myStrs) != 3 {
		return ext.EndGroups
	}

	caseId := myStrs[2]
	scan := database.GetScan(caseId)
	if scan == nil {
		return ext.EndGroups
	}

	if !scan.IsPending() {
		_, _ = replied.EditText(b, scan.ParseAsMd().ToString(), &gotgbot.EditMessageTextOpts{
			DisableWebPagePreview: true,
		})
		return ext.EndGroups
	}

	args := strongStringGo.SplitN(message.Text, 2, " ", "\n")
	var newReason string
	if len(args) > 1 {
		newReason = args[1]
	}

	scan.AgentUser = tgUser

	scan.Close(token.UserId, newReason)
	database.UpdateScan(scan)
	_, _ = replied.EditText(b, scan.ParseAsMd().ToString(), &gotgbot.EditMessageTextOpts{
		ParseMode:             sv.MarkDownV2,
		DisableWebPagePreview: true,
	})

	_, _ = replied.Delete(b)

	return ext.EndGroups
}
