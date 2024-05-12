import { defineComponent, Ref, ref, onMounted, onUnmounted, watch } from "vue";
import { debounce } from "@/utils";
import "./style.scss";

export default defineComponent({
  name: "TextTip",
  props: {
    content: {
      type: String,
      default: "",
    },
  },
  setup(props) {
    const textTipRef = ref<HTMLElement>();
    const textTipTextRef = ref<HTMLElement>();
    const disabled: Ref<boolean> = ref(true);
    let observer: ResizeObserver;

    const checkContent = () => {
      const textTipDom = textTipRef.value;
      const textTipTextDom = textTipTextRef.value;
      if (!textTipDom || !textTipTextDom) {
        return;
      }
      if (textTipTextDom.scrollWidth > textTipDom.scrollWidth) {
        disabled.value = false;
      } else {
        disabled.value = true;
      }
    };

    watch(
      () => props.content,
      () => {
        checkContent();
      }
    );

    onMounted(() => {
      checkContent();
      if (window.ResizeObserver) {
        observer = new ResizeObserver(
          debounce(() => {
            checkContent();
          }, 100)
        );
        observer.observe(textTipRef.value!);
      }
    });

    onUnmounted(() => {
      if (observer) {
        observer.disconnect();
      }
    });

    return () => {
      return (
        <div class="text-tip" ref={textTipRef} title={disabled.value ? "" : props.content}>
          <div class={"text"} ref={textTipTextRef}>
            {props.content}
          </div>
        </div>
      );
    };
  },
});
