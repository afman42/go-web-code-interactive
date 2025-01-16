import { writable } from "svelte/store";
export let langState = writable({
  value: "node",
  sampleDataLang: {
    node: "console.log('hello world')\n\n\n\n\n\n",
    php: "<?php\necho 'hello world';\n\n\n\n\n",
    go: 'package main\n\nimport "fmt"\n\nfunc main(){\n\tfmt.Println("hello world")\n}\n\n\n\n\n',
  },
});
