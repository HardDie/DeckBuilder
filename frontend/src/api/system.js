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
    return withBindingError(
      UpdateSettings({
        lang: body?.lang || '',
        enable_back_shadow: Boolean(body?.enable_back_shadow),
        card_scale: Number(body?.card_scale),
      }),
    )
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
