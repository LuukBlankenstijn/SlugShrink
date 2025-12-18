import type { Domain, FullRedirect } from '../gen/api/v1/api_pb';

export const MODAL_CONTEXT = 'modal-controls';

export type ModalControls = {
	openDomainModal: (domain?: Domain) => void;
	openRedirectModal: (redirect?: FullRedirect) => void;
	closeModal: () => void;
};
