<template>
  <span class="typing" :class="[`is-${align}`, { 'is-mono': mono }]" :style="{ fontSize: size }">
    <!-- stack：sizer 撑开最长文案的高度与宽度，避免打字过程中页面抖动 -->
    <span class="typing-stack">
      <span class="typing-sizer" aria-hidden="true">{{ longestLine }}</span>
      <span class="typing-view" aria-hidden="true">
        <span class="typing-body">{{ displayed }}</span><span v-if="caret" class="typing-caret"></span>
      </span>
    </span>
    <!-- 屏幕阅读器直接读完整文案，避免逐字播报 -->
    <span class="typing-sr">{{ readableText }}</span>
  </span>
</template>

<script setup>
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'

const props = defineProps({
  // 循环播放的文案列表
  lines: { type: Array, default: () => [] },
  // 单字打字间隔（ms），实际会做 ±30% 抖动，更接近真人手感
  typingSpeed: { type: Number, default: 95 },
  // 单字删除间隔（ms）
  deletingSpeed: { type: Number, default: 45 },
  // 整句打完后的停留时间（ms）
  holdTime: { type: Number, default: 1500 },
  // 挂载后延迟多久开始打字（ms）
  startDelay: { type: Number, default: 500 },
  // 是否循环播放；关闭时停在最后一句
  loop: { type: Boolean, default: true },
  // 是否显示光标
  caret: { type: Boolean, default: true },
  // 对齐方式：left | center
  align: { type: String, default: 'left' },
  // 是否使用等宽字体
  mono: { type: Boolean, default: true },
  // 字号，默认继承父级
  size: { type: String, default: 'inherit' }
})

// 按 Unicode 码点切分，保证 emoji / 代理对不会被拆成半个字符
const charsOf = (str) => Array.from(str || '')

const validLines = computed(() => props.lines.filter((l) => typeof l === 'string' && l.length > 0))
const longestLine = computed(() =>
  validLines.value.reduce((a, b) => (charsOf(b).length > charsOf(a).length ? b : a), '')
)
const readableText = computed(() => validLines.value.join('，'))

const displayed = ref('')
const lineIndex = ref(0)
const charIndex = ref(0)
let timer = null

const schedule = (fn, delay) => {
  clearTimeout(timer)
  timer = setTimeout(fn, delay)
}

// 打字速度抖动，避免机械感
const jitter = (speed) => Math.max(16, Math.round(speed * (0.7 + Math.random() * 0.6)))

const typeNext = () => {
  const chars = charsOf(validLines.value[lineIndex.value])
  if (charIndex.value < chars.length) {
    displayed.value += chars[charIndex.value]
    charIndex.value += 1
    schedule(typeNext, jitter(props.typingSpeed))
    return
  }
  // 打完整句：循环模式继续，否则停在最后一句
  const isLast = lineIndex.value >= validLines.value.length - 1
  if (props.loop || !isLast) schedule(startDeleting, props.holdTime)
}

const startDeleting = () => schedule(deleteNext, props.deletingSpeed)

const deleteNext = () => {
  const chars = charsOf(validLines.value[lineIndex.value])
  if (charIndex.value > 0) {
    charIndex.value -= 1
    displayed.value = chars.slice(0, charIndex.value).join('')
    schedule(deleteNext, props.deletingSpeed)
    return
  }
  lineIndex.value = (lineIndex.value + 1) % validLines.value.length
  schedule(typeNext, 250)
}

const start = () => {
  clearTimeout(timer)
  displayed.value = ''
  lineIndex.value = 0
  charIndex.value = 0
  if (!validLines.value.length) return

  // 尊重系统的「减少动态效果」设置：直接静态展示首句
  const reduceMotion =
    typeof window !== 'undefined' && window.matchMedia?.('(prefers-reduced-motion: reduce)').matches
  if (reduceMotion) {
    displayed.value = validLines.value[0]
    return
  }
  schedule(typeNext, props.startDelay)
}

onMounted(start)
// 文案异步变化时重新开始一轮
watch(() => props.lines, start, { deep: true })
onBeforeUnmount(() => clearTimeout(timer))
</script>

<style scoped>
.typing {
  display: block;
  color: var(--accent-color);
  font-weight: 500;
  line-height: 1.5;
}
.typing.is-mono {
  font-family: ui-monospace, SFMono-Regular, "Fira Code", Menlo, Consolas, "Liberation Mono", monospace;
}

/* sizer 与 view 叠在同一网格单元：sizer 撑开尺寸，view 负责显示 */
.typing-stack {
  display: grid;
}
.typing-sizer,
.typing-view {
  grid-area: 1 / 1;
  min-width: 0;
}
/* sizer 与正文使用完全相同的换行规则，保证两者行数一致、高度同步 */
.typing-sizer {
  visibility: hidden;
  white-space: pre-wrap;
  word-break: break-word;
}
.typing-view {
  display: flex;
  align-items: center;
}
.typing-body {
  flex: 0 1 auto;
  min-width: 0;
  white-space: pre-wrap;
  word-break: break-word;
}
.typing-caret {
  flex: none;
  width: 0.5em;
  height: 1.05em;
  margin-left: 0.12em;
  background: currentColor;
  animation: typing-blink 1s steps(1, end) infinite;
}

.typing.is-center .typing-stack {
  justify-items: center;
}
.typing.is-center .typing-view {
  justify-content: center;
}

.typing-sr {
  position: absolute;
  width: 1px;
  height: 1px;
  padding: 0;
  margin: -1px;
  overflow: hidden;
  clip: rect(0 0 0 0);
  white-space: nowrap;
  border: 0;
}

@keyframes typing-blink {
  50% {
    opacity: 0;
  }
}

@media (prefers-reduced-motion: reduce) {
  .typing-caret {
    animation: none;
  }
}
</style>
