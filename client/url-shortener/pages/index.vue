<template>
    <div>
        <va-app-bar absolute hide-on-scroll target="#va-app-bar-hide">
            <va-spacer />
            <va-button style="margin-right: 1rem;" icon="logout" color="text-inverted" preset="plain">
                Logout
            </va-button>
        </va-app-bar>

        <div>
            <main id="va-app-bar-hide">
                <va-sidebar class="sidebar" hoverable minimized-width="64px">
                    <va-sidebar-item v-for="item in items" :key="item.title" :active="item.active">
                        <va-sidebar-item-content>
                            <va-icon :name="item.icon" />
                            <va-sidebar-item-title>
                                {{ item.title }}
                            </va-sidebar-item-title>
                        </va-sidebar-item-content>
                    </va-sidebar-item>
                </va-sidebar>
                <div style="width: 100%; display: grid;   grid-auto-rows: 15vh;">
                    <va-form style="grid-row: 1;" tag="form" id="submition-form">
                        <va-input v-model="value" type="url" label="url"
                            :rules="[(v) => v.length > 0 || `This field is required`]" />
                        <br />
                        <va-button class="submit-button" type="submit">
                            Submit
                        </va-button>
                    </va-form>
                    <div style="grid-row: 2; grid-column: 1; top: 0; padding-left: 3vh; padding-top: 5vh; color: var(--va-text-primary) !important;"
                        class="row">
                        <va-input v-model="input" class="" placeholder="Filter..." />

                        <div style="padding-top: 2vh;">
                            <va-checkbox v-model="useCustomFilteringFn" label="Exact match" />
                            <br />
                            <va-checkbox v-model="isDebounceInput" label="Debounce input" />
                        </div>
                    </div>
                    <va-data-table style="grid-row: 3; grid-column: 1; top: 0;padding-left: 3vh; padding-top: 8vh;"
                        :columns="columns" class="data-table" striped :items="urls" :filter="filter"
                        :filter-method="customFilteringFn" @filtered="filteredCount = $event.items.length" />
                </div>
            </main>
        </div>
    </div>
</template>

<script lang="ts">
import debounce from "lodash/debounce.js";
import { defineComponent } from "vue";

export default defineComponent({
    data() {
        const columns = [
            { key: "id", sortable: true },
            { key: "upsteam_url", sortable: true, label: "upsteam url" },
            { key: "downstream_uri", sortable: true, label: "downstream uri" },
            { key: "created_at", sortable: true, label: "created at" },
            { key: "updated_at", sortable: true, label: "updated at" },
            { key: "creator" },
            { key: "updater" },
            { key: "version", sortable: true },
        ];

        const urls = [
            {
                id: 1,
                upsteam_url: "https://google.com",
                downstream_uri: "/dddddddfff",
                created_at: null,
                updated_at: null,
                creator: "admin",
                updater: "admin",
                version: 1,
            }
        ]
        const input = "";

        return {
            value: "",
            urls,
            columns,
            input,
            filter: input,
            isDebounceInput: false,
            useCustomFilteringFn: false,
            filteredCount: urls.length,
            items: [
                { title: "URLs", icon: "link", active: true },
            ],
        };
    },
    computed: {
        customFilteringFn(): any {
            return this.useCustomFilteringFn ? this.filterExact : undefined;
        },
    },

    methods: {
        filterExact(source: any): any {
            if (this.filter === "") {
                return true;
            }
            return source?.toString?.() === this.filter;
        },

        updateFilter(filter: any) {
            this.filter = filter;
        },

        debouncedUpdateFilter: debounce(function (filter) {
            this.updateFilter(filter);
        }, 600),
    },

    watch: {
        input(newValue: any) {
            if (this.isDebounceInput) {
                this.debouncedUpdateFilter(newValue);
            } else {
                this.updateFilter(newValue);
            }
        },
    },
});

</script>

<style lang="scss">
#va-app-bar-hide {
    padding-top: 0.125rem;
    overflow: auto;
    height: 100vh !important;
    width: 100% !important;
    overflow-y: scroll;
    overflow-x: hidden;
    display: flex;
    background: var(--va-background-primary);
    box-sizing: border-box;
}

#submition-form {
    justify-items: center;
    justify-content: center;
    display: block;
    padding-top: 3vh;
    padding-left: 3vh;
    width: 100%;
    align-items: center;
}

.submit-button {
    margin-top: 3vh;
    align-self: center !important;
}

.va-data-table__table-th {
    background-color: var(--va-background-element);
    color: var(--va-text-primary) !important;
}

.va-data-table__table-td {
    color: var(--va-text-primary) !important;
}
</style>
