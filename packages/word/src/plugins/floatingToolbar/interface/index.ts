import Editor from '../../../editor'

export interface IToolbarRegister {
  key?: string
  isDivider?: boolean
  render?: (container: HTMLDivElement, editor: Editor) => void
  callback?: (editor: Editor) => void
}
