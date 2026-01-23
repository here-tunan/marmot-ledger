import request from '@/api/request'

export function allTransactionAccount() {
    return request.get("/money/transactionAccount")
}

export function allTransactionAccountByFamily(familyId) {
    return request.get(`/money/transactionAccount/family/${familyId}`)
}

export function createTransactionAccount(data) {
    return request.put("/money/transactionAccount", data)
}

export function updateTransactionAccount(data) {
    return request.put("/money/transactionAccount", data)
}

export function deleteTransactionAccount(id) {
    return request.delete("/money/transactionAccount", { params: { id } })
}