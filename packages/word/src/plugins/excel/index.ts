import Editor from '../../editor'
import importExcel from './importExcel'

export default function excelPlugin(editor: Editor) {
  const command = editor.command

  command.executeImportExcel = importExcel(command)
}
