export function onStorageKeyChange(
	key: string,
	onChange: (event: StorageEvent) => void,
): () => void {
	const listener = (event: StorageEvent) => {
		if (event.key === key) {
			onChange(event);
		}
	};
	window.addEventListener("storage", listener);
	return () => window.removeEventListener("storage", listener);
}
