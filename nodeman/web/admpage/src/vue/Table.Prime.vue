<template>
  <DataTable :value="products" tableStyle="min-width: 50rem">
    <Column
      v-for="col in columns"
      :key="col.field"
      :field="col.field"
      :header="col.header"
    />
  </DataTable>
</template>

<script setup lang="ts">
import { ref } from "vue";
import DataTable from "primevue/datatable";
import Column from "primevue/column";

type ColumnDef = {
  field: string;
  header: string;
};

const columns = ref<ColumnDef[]>([
  { field: "code", header: "Code" },
  { field: "name", header: "Name" },
  { field: "category", header: "Category" },
  { field: "quantity", header: "Quantity" },
]);

const products = ref([
  {
    code: "P-1001",
    name: "Keyboard",
    category: "Electronics",
    quantity: 10,
  },
  {
    code: "P-1002",
    name: "Mouse",
    category: "Electronics",
    quantity: 25,
  },
]);
</script>
