import { createI18n } from "vue-i18n";

export const i18n = createI18n({
  legacy: false,
  locale: "en",
  fallbackLocale: "en",
  messages: {
    en: {
      table: {
        users: {
          status: {
            enabled: "Enabled",
            disabled: "Disabled",
          },
          actions: {
            enable: "Enable",
            disable: "Disable",
            delete: {
              button: "Delete",
              confirm: {
                header: "Delete user?",
                body: "Are you sure?",
                ok: "Da",
                cancel: "Net",
              },
            },
          },
          columns: {
            id: "ID",
            "display-name": "Display Name",
            "target-status": "Target Status",
            name: "Name",
            "vless-uuid": "VLESS UUID",
            actions: "Actions",
          },
        },
      },
    },
    ru: {
      table: {
        users: {
          status: {
            enabled: "Энейбл",
            disabled: "Дизейбл",
          },
          actions: {
            enable: "Enable",
            disable: "Disable",
            delete: {
              button: "Delete",
              confirm: {
                header: "Delete user?",
                body: "Are you sure?",
                ok: "Da",
                cancel: "Net",
              },
            },
          },
          id: "IIDD",
        },
      },
    },
  },
});
