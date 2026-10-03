import {
  DownloadStatus,
  GetSettings,
  GetVersion,
  Status,
  UpdateSettings,
} from '../../wailsjs/go/system/System'
import { withBindingError } from '@/api/wails'

export default {
  getSettings() {
    return withBindingError(GetSettings())
  },
  updateSettings(body) {
    return withBindingError(UpdateSettings({ lang: body?.lang || '' }))
  },
  checkStatus() {
    return withBindingError(Status())
  },
  downloadStatus() {
    return withBindingError(DownloadStatus())
  },
  getVersion() {
    return withBindingError(GetVersion())
  },
}
