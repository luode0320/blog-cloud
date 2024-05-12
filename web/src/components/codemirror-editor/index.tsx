import { defineComponent } from "vue";
import { Codemirror } from "vue-codemirror";
import { json } from "@codemirror/lang-json";
import { keymap } from "@codemirror/view";
import "./index.scss";

export default defineComponent({
  name: "CodemirrorEditor",
  props: {
    noRadius: {
      type: Boolean,
      default: false,
    },
  },
  emits: ["save"],
  setup(props, { emit }) {
    return () => {
      return (
        <Codemirror
          style={{ width: "100%" }}
          class={props.noRadius ? "no-radius" : ""}
          autofocus
          indent-with-tab
          tab-size={2}
          extensions={[
            json(),
            keymap.of([
              {
                key: "Ctrl-s",
                run: (view) => {
                  emit("save");
                  return true;
                },
              },
            ]),
          ]}
        />
      );
    };
  },
});
