<template id="main-page">
    <div>
        <va-app-bar absolute hide-on-scroll target="#va-app-bar-hide">
            <span
                style="margin-left: 2vh; font-family: var(--va-font-family); font-weight: var(--va-button-font-weight);">URL
                Shortener</span>
            <va-spacer />
            <va-button style="margin-right: 2vh;" icon="logout" color="text-inverted" preset="plain">
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
                <div style="width: 100%; display: grid; grid-auto-rows: 15vh;">
                    <va-alert center style="align-self: center; grid-row: 1; grid-column: 1; margin-left: 3vh;"
                        color="warning" outline>
                        <template #icon>
                            <va-icon name="star" color="warning" />
                        </template>
                        Let's Short LoOoOoOng Web URLs
                    </va-alert>
                    <va-button @click="$refs.modal.show()"
                        style="grid-row: 2; grid-column: 1; top: 0; margin-left: 3vh; margin-right: 1vh; margin-bottom: 0; height: 5vh;"
                        class="submit-button" type="submit">
                        New
                    </va-button>
                    <va-modal no-padding style="padding-bottom: 10%;" ref="modal" stateful :message="'hi'">
                        <template #content="{ ok, cancel }">
                            <va-form tag="form" id="submition-form" @submit.prevent="handleSubmit">
                                <va-card-content>
                                    <va-input v-model="value" type="url" label="url"
                                        :rules="[(v) => v.length > 0 || `This field is required`]" />
                                    <!-- <br /> -->
                                </va-card-content>
                                <va-card-actions>
                                    <va-button color="warning" @click="cancel">
                                        Cancel :/
                                    </va-button>
                                    <va-button color="primary" class="submit-button" type="submit" @click="ok">
                                        Ok ;)
                                    </va-button>
                                </va-card-actions>
                            </va-form>
                        </template>
                    </va-modal>

                    <div style="align-items: center; display: grid; grid-row: 3; grid-column: 1; top: 0; margin-left: 3vh; margin-top: 5vh; color: var(--va-text-primary) !important;"
                        class="row">
                        <va-input v-model="input" style="grid-column: 1;" placeholder="Filter..." />

                        <va-checkbox style="grid-column: 2; margin-left: 3vh;" v-model="useCustomFilteringFn"
                            label="Exact match" />
                        <va-checkbox style="grid-column: 3;" v-model="isDebounceInput" label="Debounce input" />
                    </div>
                    <va-data-table virtual-scroller
                        style="grid-row: 4; grid-column: 1; top: 0;margin-left: 3vh; margin-top: 2%; height: 40vh;"
                        :columns="columns" class="data-table" striped :items="urls" :filter="filter"
                        :filter-method="customFilteringFn" @filtered="filteredCount = $event.items.length">

                        <template #cell(actions)="{ rowIndex }">
                            <va-button preset="plain" icon="edit" @click="openModalToEditItemById(rowIndex)" />
                            <va-button preset="plain" icon="delete" @click="deleteItemById(rowIndex)" />
                        </template>
                    </va-data-table>
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
            { key: "actions" },
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
            },
            {
                id: 2,
                upsteam_url: "https://google.com",
                downstream_uri: "/dddddddfff",
                created_at: null,
                updated_at: null,
                creator: "admin",
                updater: "admin",
                version: 1,
            },
            {
                id: 3,
                upsteam_url: "https://google.com",
                downstream_uri: "/dddddddfff",
                created_at: null,
                updated_at: null,
                creator: "admin",
                updater: "admin",
                version: 1,
            },
            {
                id: 4,
                upsteam_url: "https://google.com",
                downstream_uri: "/dddddddfff",
                created_at: null,
                updated_at: null,
                creator: "admin",
                updater: "admin",
                version: 1,
            },
            {
                id: 5,
                upsteam_url: "https://google.com",
                downstream_uri: "/dddddddfff",
                created_at: null,
                updated_at: null,
                creator: "admin",
                updater: "admin",
                version: 1,
            },
            {
                id: 5,
                upsteam_url: "https://google.com",
                downstream_uri: "/dddddddfff",
                created_at: null,
                updated_at: null,
                creator: "admin",
                updater: "admin",
                version: 1,
            },
            {
                id: 6,
                upsteam_url: "https://google.com",
                downstream_uri: "/dddddddfff",
                created_at: null,
                updated_at: null,
                creator: "admin",
                updater: "admin",
                version: 1,
            },
            {
                id: 7,
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
        handleSubmit() {
            alert("-- form submit --");
        },
        deleteItemById(id: number) {
            alert(id);
        },
        openModalToEditItemById(id: number) {
            alert(id);
        },
        filterExact(source: any): any {
            if (this.filter === "") {
                return true;
            }
            return source?.toString?.() === this.filter;
        },

        updateFilter(filter: string) {
            this.filter = filter;
        },

        debouncedUpdateFilter: debounce(function (filter: string) {
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

.va-modal--no-padding {
    padding: 0;
    padding-bottom: 5%;
    padding-right: 8%;
}

.va-virtual-scroller {
    width: unset;
}
</style>
