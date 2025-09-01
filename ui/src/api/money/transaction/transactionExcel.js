import request from '@/api/request'

// 微信XLSX账单导入
export function postWeChatExcel(file) {
    return request.post("/money/transaction/import/wechat", file, {
        headers: {"Content-Type": "multipart/form-data"}
    })
}

// 支付宝CSV账单导入
export function postAlipayExcel(file) {
    return request.post("/money/transaction/import/alipay", file, {
        headers: {"Content-Type": "multipart/form-data"}
    })
}