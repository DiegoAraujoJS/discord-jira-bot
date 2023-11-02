package utils

var usageRecord = map[string]int{}

func ExposeUsageDetails(record string, exposeFn func(record string, records map[string]int)) {
    usageRecord[record] = usageRecord[record] + 1

    exposeFn(record, usageRecord)
}
