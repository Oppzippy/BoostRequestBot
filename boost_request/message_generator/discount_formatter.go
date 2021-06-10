package message_generator

import (
	"strings"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/oppzippy/BoostRequestBot/boost_request/repository"
	"github.com/shopspring/decimal"
)

type DiscountFormatter struct {
	roleNameProvider RoleNameProvider
	localizer        *i18n.Localizer
}

func NewDiscountFormatter(localizer *i18n.Localizer, roleNameProvider RoleNameProvider) *DiscountFormatter {
	return &DiscountFormatter{
		localizer:        localizer,
		roleNameProvider: roleNameProvider,
	}
}

func (df *DiscountFormatter) FormatDiscounts(br *repository.BoostRequest) string {
	sb := strings.Builder{}
	if len(br.RoleDiscounts) != 0 {
		for _, rd := range br.RoleDiscounts {
			sb.WriteString(df.FormatDiscount(rd))
		}
	}
	return sb.String()
}

func (df *DiscountFormatter) FormatDiscount(roleDiscount *repository.RoleDiscount) string {
	roleName := df.roleNameProvider.RoleName(roleDiscount.GuildID, roleDiscount.RoleID)
	discountPercentage := roleDiscount.Discount.Mul(decimal.NewFromInt(100))

	if roleName != "" {
		return df.localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "DiscountOnBoostTypeForRole",
				Other: "{{.DiscountPercentage}}% discount on {{.BoostType}} ({{.Role}})",
			},
			TemplateData: map[string]interface{}{
				"DiscountPercentage": discountPercentage,
				"BoostType":          roleDiscount.BoostType,
				"Role":               roleName,
			},
		})
	}

	return df.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "DiscountOnBoostType",
			Other: "{{.DiscountPercentage}}% discount on {{.BoostType}}",
		},
		TemplateData: map[string]interface{}{
			"DiscountPercentage": discountPercentage,
			"BoostType":          roleDiscount.BoostType,
		},
	})
}
