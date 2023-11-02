package utils

var usageRecord map[string]int

func ExposeUsageDetails(record string, exposeFn func(record string, records map[string]int)) {
    if val, ok := usageRecord[record]; ok {
        usageRecord[record] = val + 1
    } else {
        usageRecord[record] = 0
    }

    exposeFn(record, usageRecord)
}
