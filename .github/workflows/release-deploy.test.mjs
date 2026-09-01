import test from 'node:test'
import assert from 'node:assert/strict'
import { readFileSync } from 'node:fs'

const workflow = readFileSync(new URL('./release-deploy.yml', import.meta.url), 'utf8')

test('supports manual production releases from main with an immutable source SHA', () => {
  assert.match(workflow, /^  workflow_dispatch:$/m)
  assert.match(workflow, /github\.event_name == 'workflow_dispatch' && github\.ref == 'refs\/heads\/main'/)
  assert.match(workflow, /RELEASE_SHA: \$\{\{ github\.event_name == 'workflow_dispatch' && github\.sha \|\| github\.event\.workflow_run\.head_sha \}\}/)
  assert.equal(workflow.match(/ref: \$\{\{ env\.RELEASE_SHA \}\}/g)?.length, 2)
  assert.match(workflow, /echo "tag=sha-\$RELEASE_SHA" >> "\$GITHUB_OUTPUT"/)
})
