import { GetSettings, GetVersion, Status, UpdateSettings } from '../../wailsjs/go/system/System'
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
  getVersion() {
    return withBindingError(GetVersion())
  },
}
