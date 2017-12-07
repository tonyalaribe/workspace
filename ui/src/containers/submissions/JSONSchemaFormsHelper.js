import React from "react";
import FileWidget from "../../components/fileWidget.js";

export const log = type => console.log.bind(console, type);

export const widgets = {
	FileWidget: FileWidget
};

export function CustomFieldTemplate(props) {
	const {
		id,
		classNames,
		label,
		help,
		required,
		description,
		errors,
		children
	} = props;

	return (
		<div className={classNames + " pv2 tl"}>
			<label htmlFor={id} className="pv2 dib">
				{label}
				{required ? "*" : null}
			</label>
			{description}
			{children}
			{errors}
			{help}
		</div>
	);
}
