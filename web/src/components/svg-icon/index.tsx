import { defineComponent, computed } from "vue";
import "./style.scss";

export default defineComponent({
  name: "SvgIcon",
  props: {
    prefix: {
      type: String,
      default: "icon",
    },
    name: {
      type: String,
      required: true,
    },
    className: {
      type: String,
    },
    customStyle: {
      type: String,
    },
  },
  emits: ["click"],
  setup(props, { emit }) {
    const svgName = computed(() => {
      return `#${props.prefix}-${props.name}`;
    });
    const svgClass = computed(() => {
      return "svg-icon " + props.className;
    });
    return () => {
      return (
        <>
          <svg
            class={svgClass.value}
            aria-hidden={true}
            style={props.customStyle}
            onClick={(ev: MouseEvent) => {
              emit("click", ev);
            }}
          >
            <use xlinkHref={svgName.value}></use>
          </svg>
        </>
      );
    };
  },
});
