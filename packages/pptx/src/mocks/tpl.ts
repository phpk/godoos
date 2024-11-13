import scienceTpl from './tpl/science.json'
import summaryTpl from './tpl/summary.json'
import productTpl from './tpl/product.json'

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
  }
]