import { drawerOpenClosePause } from '../config'

export async function clearFillTextInput(pg, id, value) {
  if (value === undefined) {
    return
  }
  const input = await pg.$(`input#${id}`)
  await input.click({ clickCount: 3 })
  await input.type(value.toString())
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

/**
 * The drawer open/close animation can take a short while before elements are clickable. 
 * Call this to wait for the drawer to complete opening/closing.
 * @param {Page} pg 
 */
export async function waitForDrawerOpenClose(pg) {
  await pg.waitFor(drawerOpenClosePause)
}