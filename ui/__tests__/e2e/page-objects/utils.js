import { drawerOpenClosePause } from '../config'

export async function clearFillTextInput(pg, id, value) {
  if (value === undefined) {
    return
  }
  const input = await pg.waitForSelector(`input#${id}`)
  await input.click({ clickCount: 3 })
  await input.type(value.toString())
}

export async function setCascader(pg, id, values) {
  if (values === undefined || values.length === 0) {
    return
  }
  const input = await pg.waitForSelector(`input#${id}`)
  await input.click()
  for (let x = 0; x < values.length; x++) {
    const menuItem = await pg.waitForSelector(`li.ant-cascader-menu-item[title='${values[x]}'`)
    await menuItem.click()
  }
}

export async function setSelect(pg, id, value) {
  if (value === undefined) {
    return
  }
  const select = await pg.waitForSelector(`#${id}`)
  await select.click()
  await expect(pg).toClick('li.ant-select-dropdown-menu-item', { text: value })
}

export async function setSwitch(pg, id, setOn) {
  if (setOn === undefined) {
    return
  }
  const switchButtonsCheck = await pg.$$(`button#${id}[aria-checked='${setOn ? 'false' : 'true'}']`)
  if (switchButtonsCheck.length === 0) {
    // Already set as desired.
    return
  }
  await switchButtonsCheck[0].click()
}

export async function modalYes(pg, textToCheck) {
  if (textToCheck) {
    await expect(pg).toMatch(textToCheck)
  }
  await pg.waitForSelector('button.ant-btn-danger', { visible: true })
  await pg.click('button.ant-btn-danger')
  // The modal takes a beat to disappear, need to wait for the animation
  // else we can be in trouble.
  await waitForDrawerOpenClose(pg)
}

export async function popConfirmYes(pg, textToCheck) {
  if (textToCheck) {
    await expect(pg).toMatch(textToCheck)
  }
  // tiny wait as otherwise puppeteer thinks the button is too small to click on
  // see https://github.com/puppeteer/puppeteer/blob/6474edb9ba5ab6637361198b574dc64529eef26b/src/common/JSHandle.ts#L433
  await pg.waitFor(100)
  await expect(pg).toClick('.ant-popover-buttons .ant-btn-primary', { text: 'Yes' })
}

/**
 * The drawer open/close animation can take a short while before elements are clickable. 
 * Call this to wait for the drawer to complete opening/closing.
 * @param {Page} pg 
 */
export async function waitForDrawerOpenClose(pg) {
  await pg.waitFor(drawerOpenClosePause)
}