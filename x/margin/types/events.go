//go:build FEATURE_TOGGLE_MARGIN_CLI_ALPHA
// +build FEATURE_TOGGLE_MARGIN_CLI_ALPHA

package types

const EventOpen = "margin/mtp_open"
const EventClose = "margin/mtp_close"
const EventForceClose = "margin/mtp_force_close"
const EventInterestRateComputation = "margin/interest_rate_computation"
const EventMarginUpdateParams = "margin/update_params"
const EventRepayInsuranceFund = "margin/repay_insurance_fund"
const EventBelowRemovalThreshold = "margin/below_removal_threshold"
const EventAboveRemovalThreshold = "margin/above_removal_threshold"
const EventIncrementalPayInsuranceFund = "margin/incremental_pay_insurance_fund"

const AttributeKeyPoolInterestRate = "margin_pool_interest_rate"
const AttributeKeyMarginParams = "margin_params"
