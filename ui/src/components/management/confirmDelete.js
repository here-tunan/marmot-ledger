import { ElMessageBox } from 'element-plus'

const DEFAULT_OPTIONS = {
  type: 'warning',
  customClass: 'calm-marmot-message-box calm-marmot-delete-box',
  confirmButtonClass: 'calm-marmot-danger-confirm',
  cancelButtonClass: 'calm-marmot-soft-cancel',
}

export async function confirmDelete({ message, title, confirmText, cancelText, type }) {
  const options = {
    ...DEFAULT_OPTIONS,
    confirmButtonText: confirmText || '删除',
    cancelButtonText: cancelText || '取消',
    type: type || DEFAULT_OPTIONS.type,
  }
  return ElMessageBox.confirm(message, title, options)
}

export function isCancelError(err) {
  return err === 'cancel' || err?.message === 'cancel'
}
