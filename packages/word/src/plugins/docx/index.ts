import Editor from '../../editor'
import exportDocx from './exportDocx'
import importDocx from './importDocx'

export default function docxPlugin(editor: Editor) {
  const command = editor.command
  // 导入文档
  command.executeImportDocx = importDocx(command)
  // 导出文档
  command.executeExportDocx = exportDocx(command)
}
