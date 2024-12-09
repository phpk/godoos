import scienceTpl from './tpl/science.json'
import summaryTpl from './tpl/summary.json'
import productTpl from './tpl/product.json'
import taikongrenTpl from './tpl/taikongren.json'
import dangjianTpl from './tpl/dangjian.json'
import jianyueziTpl from './tpl/jianyuezi.json'

export const tplList = [
  {
    label: '基础',
    key: 'base',
    tpl: []
  },
  {
    label: '科技',
    key: 'science',
    tpl: scienceTpl
  },
  {
    label: '总结',
    key: 'summary',
    tpl: summaryTpl
  },
  {
    label: '产品',
    key: 'product',
    tpl: productTpl
  },
  {
    label: '太空人',
    key: 'taikongren',
    tpl: taikongrenTpl
  },
  {
    label: '党建',
    key: 'dangjian',
    tpl: dangjianTpl
  },
  {
    label: '简约紫',
    key: 'jianyuezi',
    tpl: jianyueziTpl
  }
]
