import {
  GetSettings,
  GetVersion,
  Quit,
  Status,
  UpdateSettings,
} from '../../wailsjs/go/system/System'
import { withBindingError } from '@/api/wails'

export default {
  quit() {
    return withBindingError(Quit())
  },
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
