import request from '@/api/request'

export function allTransactionCategory() {
    return request.get("/money/transactionCategory")
}

export function createTransactionCategory(data) {
    return request.put("/money/transactionCategory", data)
}

export function updateTransactionCategory(data) {
    return request.put("/money/transactionCategory", data)
}

export function deleteTransactionCategory(id) {
    return request.delete("/money/transactionCategory", { params: { id } })
}