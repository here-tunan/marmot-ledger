import { defineStore } from 'pinia';

export const useSidebarStore = defineStore('sidebar', {
	state: () => {
		return {
			collapse: false,
			showMobile: false,
		};
	},
	getters: {},
	actions: {
		handleCollapse() {
			this.collapse = !this.collapse;
		},
		openMobile() {
			this.showMobile = true;
		},
		closeMobile() {
			this.showMobile = false;
		},
		toggleMobile() {
			this.showMobile = !this.showMobile;
		}
	}
});
