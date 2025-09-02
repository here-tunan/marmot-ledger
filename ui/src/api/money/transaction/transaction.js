import request from '@/api/request'


export function putTransaction(transaction) {
    return request.put("/money/transaction", transaction)
}

export function batchPutTransaction(transactions) {
    return request.put("/money/transaction/batch", transactions)
}


export function queryTransaction(param) {
    return request.post("/money/transaction/list", param)
}

export function deleteTransaction(id) {
    return request.delete("/money/transaction", {
        params: {id: id}
    })
}

export function transactionAggregation(startTime, endTime, aggregationType) {
    return request.get("/money/transaction/aggregation", {
        params: {
            startTime: startTime,
            endTime: endTime,
            aggregationType: aggregationType
        }
    })
}

export function transactionAnalysis(param) {
    return request.post("/money/transaction/analysis", param)
}

export function checkAccountUsage(accountId) {
    return request.get("/money/transactionAccount/usage-check", {
        params: { accountId: accountId }
    })
}

export function checkCategoryUsage(categoryId) {
    return request.get("/money/transactionCategory/usage-check", {
        params: { categoryId: categoryId }
    })
}