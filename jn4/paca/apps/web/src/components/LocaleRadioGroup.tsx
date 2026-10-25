import {
	DropdownMenuRadioGroup,
	DropdownMenuRadioItem,
} from "@/components/ui/dropdown-menu";
import { SUPPORTED_LOCALES } from "@/hooks/use-locale";
import type { SupportedLanguage } from "@/i18n";

interface LocaleRadioGroupProps {
	value: SupportedLanguage;
	onValueChange: (value: SupportedLanguage) => void;
}

export function LocaleRadioGroup({
	value,
	onValueChange,
}: LocaleRadioGroupProps) {
	return (
		<DropdownMenuRadioGroup
			value={value}
			onValueChange={(next) => onValueChange(next as SupportedLanguage)}
		>
			{SUPPORTED_LOCALES.map((option) => (
				<DropdownMenuRadioItem key={option.code} value={option.code}>
					{option.nativeLabel}
				</DropdownMenuRadioItem>
			))}
		</DropdownMenuRadioGroup>
	);
}
