import REGION_DATA from 'china-area-data'
import { TextToCode } from 'element-china-area-data'

/**
 * 获取全国相应行政区域的数据
 * @param areaCode 相应区域的区域码，默认86，即中国大陆国际区号，返回全国所有省份数据。传相应省的号码即获得相应省所有市的信息，依此类推。
 * @returns value为相应区域号，label为相应区域中文名
 */
function getAreaOptions(ast = '86') {
  const provinceObj = REGION_DATA[ast]
  const regionData = []

  for (const prop in provinceObj) {
    regionData.push({
      value: provinceObj[prop],
      label: provinceObj[prop]
    })
  }

  return regionData
}

export { TextToCode, getAreaOptions }