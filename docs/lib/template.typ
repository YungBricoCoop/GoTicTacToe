// SPDX-FileCopyrightText: 2025 Haute école d'ingénierie et d'architecture de Fribourg
//
// SPDX-License-Identifier: CC0-1.0

#import "@preview/chic-hdr:0.5.0": *
#import "@preview/datify:0.1.4": custom-date-format, month-name
#import "@preview/hydra:0.6.1": hydra

#let heading_text(..content) = text(font: "Aptos", ..content)

#let cahier_cover(
  logo_path: "img/logo-heiafr.svg",
  title_line_1: "",
  doc_type: "Cahier des charges",
  cover_image_path: none,
  stream: "",
  student: "",
  supervisors: "",
  assistants: "",
  date: "",
  version: "1.0",
) = {
  let info_cols = (28mm, 90mm)
  let info_row_gutter = 6mm
  let cover_image_width = 60mm
  set page(paper: "a4", margin: 25mm)
  set text(font: "Aptos", size: 11pt)
  set align(center)

  image(logo_path, width: 95mm)

  heading_text(size: 20pt, weight: "semibold", title_line_1)
  v(2mm)

  heading_text(size: 13pt, weight: "semibold", doc_type)
  v(18mm)

  if cover_image_path != none {
    image(cover_image_path, width: cover_image_width)
    v(22mm)
  } else {
    v(22mm)
  }

  let cols = info_cols
  let col_gutter = 14mm
  let grid_w = cols.at(0) + cols.at(1) + col_gutter
  let left_pad = 40mm
  let label = it => heading_text(weight: "semibold", it)

  set align(left)
  pad(left: left_pad, block(width: grid_w, {
    grid(
      columns: cols,
      column-gutter: col_gutter,
      row-gutter: info_row_gutter,

      [#label("Filière")], [#heading_text(stream)],
      [#label("Étudiant")], [#heading_text(student)],
      [#label("Superviseur")], [#heading_text(supervisors)],
      ..(if assistants != "" { ([#label("Assistants")], [#heading_text(assistants)]) } else { () }),
      [#label("Date")], [#heading_text(date)],
      [#label("Version")], [#heading_text(version)],
    )
  }))
}


//
// Report
//
#let report(
  title: none,
  short_title: none,
  subtitle: none,
  theme_color: rgb(195, 40, 35, 255),
  isc_icon: "neg",
  cover_image_path: none,
  type: none,
  year: none,
  location: [Fribourg],
  profile: none,
  authors: (),
  versions: (),
  supervisors: (),
  supervisors_label: [Superviseurs],
  assistants: (),
  assistants_label: [Assistants],
  clients: (),
  clients_label: [Mandants],
  doc,
) = {
  let date_text = if versions.len() > 0 {
    let d = versions.last().date
    let j = d.day()
    if j == 1 {
      j = [#j#super[er]]
    }
    [#j #month-name(d.month(), "fr") #d.year()]
  } else {
    [N/A]
  }
  set text(font: "Aptos", size: 11pt)
  set page(
    paper: "a4",
    header: context hydra(1),
  )
  set heading(numbering: "1.1.1")
  set par(justify: true)

  // Heading formating for all levels
  show heading: it => {
    if it.level <= 1 {
      block(inset: (y: 10pt), heading_text(it))
    } else {
      block(inset: (y: 10pt), heading_text(it))
    }
  }
  cahier_cover(
    title_line_1: title,
    doc_type: type,
    cover_image_path: cover_image_path,
    stream: "ISC / Orientation ID",
    student: authors.map(a => a.firstname + " " + a.lastname).join(", "),
    supervisors: supervisors.join(", "),
    assistants: if assistants.len() > 0 { assistants.join(", ") } else { "" },
    date: date_text,
    version: if versions.len() > 0 { versions.last().version } else { "N/A" },
  )

  let authors_list = authors.map(a => a.firstname + " " + a.lastname).join(", ")

  // Table of content
  outline()
  pagebreak()

  show: chic.with(
    chic-footer(
      left-side: text(smallcaps((authors_list))),
      right-side: text(context (counter(page).display())),
    ),
    chic-header(
      left-side: text(emph(chic-heading-name(fill: true))),
      right-side: text(smallcaps(short_title)),
    ),
    chic-separator(.2mm, on: "footer"),
    chic-separator(.2mm, on: "header"),

    chic-offset(7mm),
    chic-height(28mm),
  )


  // Main content
  set page(numbering: "1")
  counter(page).update(1)

  doc
}
