<template id="main-page">
    <div>
        <va-app-bar absolute hide-on-scroll target="#va-app-bar-hide">
            <span
                style="margin-left: 2vh; font-family: var(--va-font-family); font-weight: var(--va-button-font-weight);">URL
                Shortener</span>
            <va-spacer />
            <va-button style="margin-right: 2vh;" icon="logout" color="text-inverted" preset="plain" @click="handleLogout">
                Logout
            </va-button>
        </va-app-bar>

        <div>
            <main id="va-app-bar-hide">
                <va-sidebar class="sidebar" hoverable minimized-width="64px">
                    <va-sidebar-item v-for="item in sidebarItems" :key="item.title" :active="item.active">
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
                    <va-button @click="($refs.modal as any).show()" icon="add"
                        style="grid-row: 2; grid-column: 1; top: 0; margin-left: 3vh; margin-right: 1vh; margin-bottom: 0; height: 5vh;"
                        class="submit-button" type="submit">
                        New
                    </va-button>
                    <va-modal no-padding style="padding-bottom: 10%;" ref="modal" stateful>
                        <template #content="{ ok, cancel }">
                            <va-form tag="form" id="submition-form" @submit.prevent="handleSubmit">
                                <va-card-content>
                                    <va-input v-model="value" type="url" label="url"
                                        :rules="[(v: string) => v.length > 0 || `This field is required`]" />
                                </va-card-content>
                                <va-card-actions>
                                    <va-button icon="backspace" color="warning" @click="cancel">
                                        Cancel
                                    </va-button>
                                    <va-button icon="done" color="primary" class="submit-button" type="submit" @click="ok">
                                        OK
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

                        <va-input style="grid-row: 2;grid-column: 1; margin-left: 0; margin-top: 1vh;"
                            v-model.number="perPage" class="flex flex-col mb-2 md3" type="number" placeholder="Items..."
                            label="Items per page" @change="fetchData()" />

                        <va-input style="grid-row: 2;grid-column: 2; margin-left: 1vh; margin-top: 1vh;"
                            v-model.number="currentPage" class="flex flex-col mb-2 md3" type="number" placeholder="Page..."
                            label="Current page" @change="fetchData()" />
                    </div>
                    <va-data-table virtual-scroller :loading="loading" :per-page="perPage" :current-page="currentPage"
                        style="grid-row: 4; grid-column: 1; top: 0;margin-left: 3vh; margin-top: 2%; height: 40vh;"
                        :columns="columns" class="data-table" striped :items="urls" :filter="filter"
                        :filter-method="customFilteringFn" @filtered="filteredCount = $event.items.length">

                        <template #cell(actions)="{ rowIndex }">
                            <va-button preset="plain" icon="edit" @click="openModalToEditItemById(rowIndex)" />
                            <va-button preset="plain" icon="delete" @click="deleteItemById(rowIndex)" />
                        </template>
                        <template #bodyAppend>
                            <tr>
                                <td>
                                    <div>
                                        <va-pagination v-model="currentPage" input :pages="pages" @change="fetchData()" />
                                    </div>
                                </td>
                            </tr>
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
import { useToast, VaModal } from 'vuestic-ui';
const { init: newToast } = useToast();

interface url {
    id: number,
    upstream_url: string,
    downstream_uri: string,
    created_at: string | null,
    updated_at: string | null,
    creator: string,
    updater: string,
    version: number,
};

interface sidebarItem {
    title: string,
    icon: string,
    active: boolean,
};

interface tokenInfo {
    expire_time: number,
    token: string,
    token_type: string,
};

interface newURLReq {
    upstream_url: string,
};

interface statusInfo {
    first_run: boolean,
    number_of_urls: number,
}

export default defineComponent({
    async setup() {
        const columns = [
            { key: "id", sortable: true },
            { key: "upstream_url", sortable: true, label: "upstream url" },
            { key: "downstream_uri", sortable: true, label: "downstream uri" },
            { key: "created_at", sortable: true, label: "created at" },
            { key: "updated_at", sortable: true, label: "updated at" },
            { key: "creator" },
            { key: "updater" },
            { key: "version", sortable: true },
            { key: "actions" },
        ];

        const status: statusInfo = {
            first_run: false,
            number_of_urls: 0,
        };

        return {
            value: "",
            urls: [],
            columns,
            input: "",
            filter: "",
            isDebounceInput: false,
            useCustomFilteringFn: false,
            loading: true,
            filteredCount: 0,
            perPage: 5,
            currentPage: 1,
            pages: 1,
            status: status,
            sidebarItems: [
                { title: "URLs", icon: "link", active: true },
            ] as sidebarItem[],
        };
    },
    async mounted() {
        console.log("mounted");
        await this.fetchData();
        this.loading = false;
        console.log("fetch");
        console.table(this.urls);
        console.table(this.status);
        if (this.status.first_run === true) {
            newToast("TODO: Register on first run");
        }
    },
    computed: {
        filteredUrls() {
            if (!this.urls) {
                return [];
            }
            return this.urls.filter((url: url) =>
                url.upstream_url.toLowerCase().includes(this.filter.toLowerCase())
            );
        },
        filteredCount(): number {
            return this.filteredUrls.length;
        },
        customFilteringFn(): any {
            return this.useCustomFilteringFn ? this.filterExact : undefined;
        },
    },

    methods: {
        getTokenInfo(): tokenInfo {
            const rawTokenInfo = localStorage.getItem('TokenInfo')
            if (rawTokenInfo === null) {
                this.$router.push("/auth/login");
                newToast({
                    message: "Login.!",
                    color: 'danger',
                });
            }
            return JSON.parse(rawTokenInfo as string);
        },
        async fetchData() {
            const tokenInfo = this.getTokenInfo();
            const config = useRuntimeConfig();
            await useFetch(config.public.baseURL + "/api/urls/", {
                method: "GET",
                headers: {
                    Authorization: tokenInfo.token_type + " " + tokenInfo.token
                },
                params: {
                    limit: this.perPage,
                    offset: this.currentPage - 1,
                },
            }).catch((error: any) => {
                console.error(error);
            }).then((data) => {
                this.urls = (data as any).data.value;
            });

            await useFetch(config.public.baseURL + "/api/status/", {
                method: "GET",
            }).catch((error: any) => {
                console.error(error);
            }).then((data) => {
                this.status = ((data as any).data.value as statusInfo);
                this.pages = this.perPage && this.perPage !== 0 ? Math.ceil((data as any).data.number_of_urls / this.perPage) : ((data as any).data.number_of_urls);
            }).finally(() => {
                this.loading = false
            });
        },
        async handleSubmit() {
            const tokenInfo = this.getTokenInfo();
            const config = useRuntimeConfig();
            const reqBody: newURLReq = {
                upstream_url: this.value,
            };
            await useFetch(config.public.baseURL + "/api/user/urls/", {
                method: "POST",
                headers: {
                    Authorization: tokenInfo.token_type + " " + tokenInfo.token
                },
                body: reqBody,
            }).catch((error: any) => {
                console.error(error);
            }).then((data) => {
                this.urls = (data as any).data.value;
            });
        },
        handleLogout() {
            localStorage.removeItem('TokenInfo');
            this.$router.push("/auth/login");
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

.va-inner-loading--active {
    margin-top: inherit;
}
</style>
