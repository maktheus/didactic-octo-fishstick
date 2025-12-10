*************
greek-fontenc
*************
Greek font encoding definition files
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

:Version: 2.6 (changelog_)

:Copyright: © 2010 -- 2023 Günter Milde <milde@users.sf.net>
:Licence:   This work may be distributed and/or modified under the
            conditions of the `LaTeX Project Public License`_, either
            version 1.3 of this license or any later version.

:Homepage:  https://codeberg.org/milde/greek-tex

:Latest Release: https://ctan.org/pkg/greek-fontenc

:Abstract: The `greek-fontenc` bundle provides encoding definition files
           for `Greek text font encodings`_ that define LICR [#]_ macros
           for characters from the Greek script

           Included are also the LaTeX packages textalpha_ and alphabeta_.

.. [#] `LaTeX internal character representation` (LICR) macros can
   serve as a human readable 7-bit-ASCII character encoding that
   works unaltered under both, 8-bit TeX and XeTeX/LuaTeX.
   Use cases are macro definitions and generated text.

.. contents::


TeX files and packages
======================

Packages
--------

.. _textalpha:

`<textalpha.sty>`_: `Greek symbols in text <textalpha.sty.html>`_
  Use ``\textalpha`` ... ``\textOmega`` or Greek literal characters
  independent of font encoding and TeX engine.

  Documentation: textalpha-doc.pdf_

  .. _alphabeta:

`<alphabeta.sty>`_: `Greek symbols in text and math <alphabeta.sty.html>`_
  Use ``\alpha`` ... ``\Omega`` independent of text/math mode,
  font encoding, and TeX engine.

  Documentation: alphabeta-doc.pdf_

Font encoding definitions
-------------------------

.. _LGR font encoding definition file:

`<greek-fontenc.def>`_
  `Common Greek font encoding definitions <greek-fontenc.def.html>`_

`<lgrenc.def>`_
  `LGR Greek font encoding definitions. <lgrenc.def.html>`_

  .. _tuenc-greek:

`<tuenc-greek.def>`_
  `Extended Greek definitions for the TU font encoding <tuenc-greek.def.html>`_

`<puenc-greek.def>`_
  `Extended Greek definitions for PDF strings <puenc-greek.def.html>`_

Auxiliary files
---------------

`<greek-euenc.def>`_
  Backwards compatibility file loading tuenc-greek.def_.
`<lgr2licr.lua>`_
  `LGR Transcription to Greek LICR transformation <lgr2licr.lua.html>`_.
  Provisional.

The source files can be converted with PyLit_ to reStructuredText_ and
with Docutils_ to the HTML documentation.


Usage examples and test documents
=================================

`<char-list.tex>`_ : `<char-list.pdf>`_, `<char-list-tu.pdf>`_
  List of Greek characters supported by `greek-fontenc`.
  Compares input variants and tests chase changing.

`<char-list-alphabeta.tex>`_ : `<char-list-alphabeta.pdf>`_, `<char-list-alphabeta-tu.pdf>`_
  List/test of Greek characters supported by `alphabeta`.


`<hyperref-with-greek.tex>`_ : `<hyperref-with-greek.pdf>`_
  Hyperref_ test and usage example.

`<test-lgrenc.tex>`_ : `<test-lgrenc.pdf>`_
  LGR test and usage example.

`<test-tuenc-greek.tex>`_ : `<test-tuenc-greek.pdf>`_
  TU test and usage example.

`<test-luainputenc.tex>`_ : `<test-luainputenc.pdf>`_
  Test LICRs with LuaTeX in 8-bit compatibility mode (with luainputenc_).


Download and Installation
=========================

The simplest way is to install this package from your distribution using
its installation manager.

Alternatively:

* Download the latest `release`_ from the package's `CTAN page`_ or a
  snapshot_ of the `greek-tex`_ repository.

* Unpack the source archive to a temporary location.

* Copy/Move/Link files ending in ``.def`` or ``.sty`` to a suitable place in
  the TeX search path.

.. _release:
    https://mirrors.ctan.org/language/greek/greek-fontenc.zip
.. _CTAN page: https://www.ctan.org/pkg/greek-fontenc
.. _greek-tex: https://codeberg.org/milde/greek-tex/
.. _snapshot: https://codeberg.org/milde/greek-tex/archive/master.zip


Conflicts
=========

The arabi_ package provides the Babel ``arabic`` option which loads
``arabicfnt.sty`` for font setup. This package overwrites the LICR macros
``\omega`` and ``\textomega`` with font selecting commands.  See the report
for Debian `bug 858987`_ for details and the `arabi workaround`_ below.

.. _bug 858987: https://bugs.debian.org/cgi-bin/bugreport.cgi?bug=858987

Usage
=====

There are several alternatives to set up the support for a Greek font
encoding provided by this bundle, e.g.:

Babel:
  Use the ``greek`` option with Babel_::

     \usepackage[greek]{babel}

  This automatically loads ``lgrenc.def`` with 8-bit TeX and
  ``tuenc-greek.def`` with XeTeX/LuaTeX and provides localized auto-strings,
  hyphenation and other localizations (see babel-greek_).

  Babel can be used together with textalpha_ or alphabeta_.

textalpha_:
  Ensure support for Greek characters in text mode::

     \usepackage{textalpha}

  eventually with the normalize-symbols_ option to handle `symbol variants`_
  and/or the keep-semicolon_ option to use the `semicolon as erotimatiko`_
  also in LGR ::

     \usepackage[normalize-symbols,keep-semicolon]{textalpha}

  This sets up LICR macros for Greek text characters under both, 8-bit TeX
  and Xe-/LuaTeX.
  For details see `<textalpha-doc.tex>`_ and `<textalpha-doc.pdf>`_ (8-bit
  TeX) as well as `<test-tuenc-greek.tex>`_ and `<test-tuenc-greek.pdf>`_
  (XeTeX/LuaTeX).

  .. _normalize-symbols: textalpha.sty.html#normalize-symbols
  .. _keep-semicolon: textalpha.sty.html#keep-semicolon
  .. _semicolon as erotimatiko: textalpha.sty.html#semicolon-as-erotimatiko

alphabeta_:
  To use the short macro names (``\alpha`` ... ``\Omega``) known from math
  mode in both, text and math mode, write ::

     \usepackage{alphabeta}

  For details see `<alphabeta-doc.tex>`_ and `<alphabeta-doc.pdf>`_.

fontenc:
  Declare LGR via fontenc_. For example, specify T1 (8-bit
  Latin) as default font encoding and LGR for Greek with ::

     \usepackage[LGR,T1]{fontenc}

  Note that without textalpha_ or alphabeta_, Greek text macros work
  only if the current font encoding supports Greek. See [fntguide]_ for
  details and `<test-lgrenc.tex>`_ for an example.

  It is possible to use 8-bit Greek text fonts in the LGR TeX font encoding
  also with XeTeX/LuaTeX, if the fontenc_ package is loaded before
  Babel, textalpha_, or alphabeta_, e.g. ::

    \usepackage[LGR]{fontenc}
    \usepackage{fontspec}
    \setmainfont{Linux Libertine O} % Latin Modern does not support Greek
    \setsansfont{Linux Biolinum O}
    \usepackage{textalpha}

  See `<test-tuenc-greek.tex>`_, `<test-tuenc-greek.pdf>`_ and
  `<test-lgrenc.tex>`_, `<test-lgrenc.pdf>`_.

.. _arabi workaround:

To work around the conflict with arabi_, it may suffice to ensure ``greek``
is loaded after ``arabic``::

    \usepackage[arabic,greek,english]{babel}

More secure is an explicit reverse-definition, e.g. ::

    % save original \omega
    \let\mathomega\omega

    \usepackage[utf8]{inputenc}
    \usepackage[LAE,LGR,T1]{fontenc}
    \usepackage[arabic,greek,english]{babel}

    % fix arabtex:
    \DeclareTextSymbol{\textomega}{LGR}{119}
    \renewcommand{\omega}{\mathomega}


Greek text font encodings
=========================

LGR
---

The LGR font encoding is the de-facto standard for typesetting Greek with
8-bit LaTeX. `greek-fontenc` provides a comprehensive `LGR font
encoding definition file`_.

Fonts in this encoding include the `CB fonts`_ (matching CM), grtimes_
(Greek Times), Kerkis_ (matching URW Bookman), DejaVu_, `Libertine GC`_, and
the `GFS fonts`_. Setup of these fonts as Greek variant to
matching Latin fonts is facilitated by the
``\DeclareFontfamilySubstitution`` command added to the
LaTeX kernel in the 2020-02 release [ltnews31]_.

The LGR font encoding allows to access Greek characters via an ASCII
transliteration. This enables simple input with a Latin keyboard.
Characters with diacritics can be selected by ligature definitions in the
font (see [greek-usage]_, [teubner-doc]_, [cbfonts]_).

A major drawback of the transliteration is the fact, that you cannot
access Latin letters if LGR is the active font encoding (e.g. in
documents or parts of documents given the `Babel` language ``greek`` or
``polutionikogreek``). This means that for every Latin-written word or
acronym an explicit language-switch is required. This problem can be
circumvented using Unicode fonts (font encoding TU_) with XeTeX or
LuaTeX.

TU
--

Standard Unicode font encoding for XeTeX and LuaTeX loaded by fontspec_
(since v2.5a) rsp. the LaTeX kernel since 2017/01/01 [ltnews26]_. [#]_
`greek-fontenc` adds support for the Greek script (see tuenc-greek_).

Xe/LuaTeX works with any system-wide installed `OpenType font`_. Suitable
fonts supporting Greek include `CM Unicode`_, `Deja Vu`_, `EB Garamond`_,
the `GFS fonts`_, `Libertine OTF`_, `Libertinus`_, `Old Standard`_,
Tempora_, and `UM Typewriter`_ (all available on CTAN) but also many commercial
fonts. Unfortunately, the fontspec_ default, `Latin Modern`_ misses most
Greek characters.

LuaTeX does not apply the NFC normalization by default. This leads to
sub-optimal placing of some diacritics, especially the sub-iota (becoming
unintelligible in combination with small letter eta). This issue can be fixed
specifiying the "Harfbuzz" renderer when loading fonts with fontspec,
e.g. ::

   \setmainfont[Renderer=Harfbuzz]{FreeSerif}

.. [#] The legacy Unicode font encodings EU1 and EU2 for XeTeX and LuaTeX
   respectively were superseded by TU in the 2017 fontspec release.

PU
--

The package hyperref_ defines the PU font encoding for use in PDF strings
(ToC, bookmarks). `greek-fontenc` adds support for Greek LICRs
(see `<hyperref-with-greek.tex>`_, `<hyperref-with-greek.pdf>`_).

----------------------------------------------------------------------------

The following two encodings are not supported by `greek-fontenc`:

LGI
---

The ‘Ibycus’ fonts from the package ibygrk_ implement an alternative
transliteration scheme (also explained in [babel-patch]_).
It is currently not supported by `greek-fontenc`.

The font encoding file ``lgienc.def`` from ibycus-babel_ provides a basic
setup (without any LICR macros or composite definitions).

T7
--

The [encguide]_ reserves the name T7 for a Greek `standard font encoding`.
However, up to now, there is no agreement on an implementation because the
restrictions for general text encodings are too severe for typesetting
polytonic Greek.


Greek LICR macro names
======================

.. note::   The LICR macro names for Greek symbols are chosen pending
            endorsement by the TeX community and related packages.

            Names for archaic characters, accents/diacritics, and
            punctuation may change in future versions.

This bundle provides LaTeX internal character representations (LICR macros)
for Greek letters and diacritics. Macro names were selected based on the
following considerations:

Letters and symbols
-------------------

* The fntguide_ (section 6.4 Naming conventions) recommends:

     Where possible, text symbols should be named as ``\text`` followed
     by the **Adobe glyph name**: for example ``\textonequarter`` or
     ``\textsterling``. Similarly, math symbols should be named as
     ``\math`` followed by the glyph name, for example
     ``\mathonequarter`` or ``\mathsterling``.

  Problem:
     The `Adobe Glyph List For New Fonts`_ has names for many glyphs in the
     `Greek and Coptic` Unicode block, but not for `Greek extended`. The
     `Adobe Glyph List`_ (for existing fonts) lists additional glyph names
     used in older fonts.  However, these are not intended for active use.

* If there exists a **math-mode macro** for a symbol, the corresponding text
  macro could be formed by prepending ``text``.

  Example:
     The glyph name for the GREEK SMALL LETTER FINAL SIGMA is ``sigma1``,
     the corresponding math-macro is ``\varsigma``. The text symbol is
     made available as ``\textvarsigma``.

  Problem:
     `Symbol variants`_ (see below).

* The `Unicode names list`_ provides standardized descriptive names for all
  Unicode characters that use only capital letters of the Latin alphabet.
  While not suited for direct use in LICR macros, they can be
  converted to LICR macro names via a defined set of transformation rules.

  Example:
    ``\textfinalsigma`` is a descriptive alias for
    GREEK SMALL LETTER FINAL SIGMA derived via the rules:

    * drop "LETTER" if the name remains unique,
    * drop "GREEK" if the name remains unique,
    * use capitalized name for capital letters, lowercase for "SMALL" letters
      and drop "SMALL",
    * concatenate

* Omit the "text" prefix for macros that do not have a math counterpart?

  Pro:
    + Simpler,
    + ease of use (less typing, better readability of source text),
    + many established text macro names without "text",
    + ``text`` prefix does **not** mark a macro as encoding-specific or
      "inserting a glyph". There are e.g. font-changing macros (``\textbf``,
      ``\textit``) and encoding-changing macros (``\textcyr``).
    + There are examples of encoding-specific macros
      without the ``text``-prefix, especially for letters, see encguide_.

  Contra:
    - Less consistent,
    - possible name clashes
    - ``text`` prefix marks a macro as confined to text (as opposed to math)
      mode,

  The font encoding definition files use the ``text`` prefix for symbols.
  Aliases (short forms, compatibility defs, etc.) are defined in
  additional packages (e.g. alphabeta.sty_ and teubner_)


Accent macros
-------------

* standard accent macros (``\DeclareTextAccent`` definitions in
  ``latex/base/...``) are one-character macros (``\' \" ... \u \v ...``) .

* ``tipa.sty``, xunicode_, and ucs_ use the "text" prefix also for accent
  macros.

  However, the `Adobe Glyph List For New Fonts`_ maps, e.g., "tonos" and
  "dieresistonos" to the spacing characters GREEK TONOS rsp. GREEK DIALYTIKA
  TONOS, hence ``\texttonos`` and ``\textdieresistonos`` should be spacing
  characters.

* textcomp (ts1enc.def) defines ``\capital...`` accents (i.e. without
  ``text`` prefix).

Currently, `greek-fontenc` uses for diacritics:

- Greek names like in Unicode, and ``ucsencs.def``, and

- the prefix ``\acc`` to distinguish the macros as `TextAaccent` and
  reduce the risk of name clashes with spacing characters.

Aliases to the "symbol macros" ``\~ \' \` \" \"' \"` ...`` are
provided. With textalpha_ or alphabeta_ also ``\<`` and ``\>`` for
``\accdasia`` and ``\accpsili``.


Symbol variants
---------------

Mathematical notation distinguishes variant shapes for beta (β|ϐ),
theta (θ|ϑ), phi (φ|ϕ), pi (π|ϖ), kappa (κ|ϰ), rho (ρ|ϱ), Theta (Θ|ϴ),
and epsilon (ε|ϵ).

The variations have no syntactic meaning in Greek text and Greek text
fonts use the shape variants indiscriminately (cf. `glyph variants`__).
The variant shapes are not given separate code-points in the LGR_ text
font encoding.

In mathematical mode, TeX supports the alternative glyph variants with
``\var<lettername>`` macros (variant macros for ϴ, ϐ, and ϰ require
additional packages).

Unicode defines separate code points for the symbol variants for use in
mathematical context. [#]_ Unfortunately, the mapping between Unicode's
letter/symbol distinction and "normal"/variant in TeX is inconsistent.

`greek-fontenc` provides ``\text<lettername>symbol`` LICR macros for the
Greek symbol characters:

* With Unicode fonts, the macros select the GREEK <lettername> SYMBOL``.

* With LGR encoded fonts, they report an error by default.

  With the ``normalize-symbols`` option of textalpha_ and alphabeta_,
  they are mapped to the corresponding letter (loosing the distinction
  between the shape variants).

The `alphabeta`_ package provides ``\<lettername>``, ``\var<lettername>``,
and ``\<lettername>symbol`` in both, text and math mode (cf. Table 1 in
`<alphabeta-doc-tu.pdf>`_).


.. [#] However, they are sometimes also used in place of the
   corresponding letter characters in Unicode-encoded text.

__ http://en.wikipedia.org/wiki/Greek_alphabet#Glyph_variants


Changelog
=========

0.9 (2013-07-03)
    - ``greek-fontenc.def`` "outsourced" from ``lgrxenc.def``
    - experimental LICRs for XeTeX/LuaTeX.
0.9.1 (2013-07-18)
    - Bugfix: wrong breathings psilioxia -> dasiaoxia.
0.9.2 (2013-07-19)
    - Bugfix: Disable composite defs starting with char macro,
    - Fix "hiatus" handling.
0.9.3 (2013-07-24)
    - Fix path for ``\input`` of ``greek-fontenc.def``.
0.9.4 (2013-09-10)
    - ``greek-fontenc.sty``: Greek text font encoding setup package.
    - remove ``xunicode-greek.sty``.
0.10 (2013-09-13)
    - textalpha_ and alphabeta_ moved here from lgrx and updated to work
      with XeTeX/LuaTeX.
    - ``greek-fontenc.sty`` removed (obsoleted by textalpha_).
0.10.1 (2013-10-01)
    - Bugfix in ``greek-euenc.def`` and ``alphabeta-euenc.def``.
0.11 (2013-11-28)
    - Compatibility with Xe/LuaTeX in 8-bit mode.
    - ``\greekscript`` *TextCommand* (cf. [encguide]_).
0.11.1 (2013-12-01)
    - Fix identification of ``greek-euenc.def``.
0.11.2 (2014-09-04)
    - Documentation update, remove duplicate code.
0.12 (2014-12-25)
    - Fix auxiliary macro names in textalpha_.
    - Conservative naming: move definition of ``\<`` and ``\>`` from
      ``greek-fontenc.def`` to ``textalpha.sty`` (Bugreport David Kastrup).
0.13 (2015-09-04)
    - Support for `symbol variants`_,
    - ``keep-semicolon`` option in textalpha_,
    - ``\lccode``/``\uccode`` corrections for Unicode
      (from Apostolos Syropoulos’ xgreek_) in greek-euenc.
    - Do not convert ``\ypogegrammeni`` to ``\prosgegrammeni``
      with ``\MakeUppercase``.
0.13.1 (2015-12-07)
    - Fix `rho with dasia bug`__ in lgrenc.def (Linus Romer).
0.13.2 (2016-02-05)
    - Support for standard Unicode text font encoding "TU"
      (new in fontspec v2.5a).
0.13.3 (2019-07-10)
    - Drop error font declaration (cf. `ltxbugs 4399`_).
0.13.4 (2019-07-11)
    - "Lowercase" ``\prosgegrammeni`` -> ``\ypogegrammeni``
      but not vice versa.
0.14 (2020-02-28)
    - Rename ``greek-euenc`` to ``tuenc-greek``.
    - Use ``\UTFencoding`` instead of ``\LastDeclaredEncoding``.
1.0 (2020-09-25)
    - Bugfix in textalpha_: Let ``\greekscript`` set ``\encodingdefault``.
    - ``\textKoppa`` as alias for ``\textkoppa`` in LGR.
2.0 (2020-10-30)
    - Move common alias definitions to ``greek-fontenc.def``.
    - textalpha_ loads TU with Xe/LuaTeX by default and provides
      ``\textmicro`` and LICR macros for archaic symbols from the
      "Greek and Coptic" Unicode block.
    - Use ``\UnicodeEncodingName`` (by the LaTeX kernel) instead of
      ``\UTFencname`` for the Unicode font encoding name.
    - Replace utf8 literals in ``tuenc-greek.def``.
    - New file ``puenc-greek.def``: setup for PU encoding defined by
      hyperref_ for PDF strings.
    - Don't use ``\textcompwordmark`` as base in accent commands.
2.1 (2022-06-14)
    - Support the correct spelling ``\guillemet…`` for « and ».
      See https://github.com/latex3/latex2e/issues/65
2.2 (2023-02-28)
    - Use correct glyph for ``\textanoteleia`` (middle dot) in LGR.
    - Test and add composite commands for combinations that are not
      converted to pre-composed characters.
    - Don't use ``\makeatother`` in ``\AtBeginDocument``.
    - Skip ``\uccode`` fixes when ignored by ``\MakeUppercase``.
    - Various small fixes and documentation update.
2.2.1 (2023-03-08)
    - Fix broken links in README.md.
    - ``@uclclist`` entry for ``\accoxia``, prevent
      downcasing ``\textStigma`` to ``\textvarstigma``.
2.2.2 (2023-03-17)
    - Don't map active ``;`` to ``\textsemicolon`` in math mode.
2.3 (2023-06-01)
    - Fix Unicode errors with pdfLaTeX and "new" (2023) ``\MakeUppercase``.
    - Upcase symbol variants also if input as LICR.
2.4 (2023-08-15)
    - Fixes for the 2022 implementation of ``\MakeUppercase``.
    - textalpha_: Map character 00B5 MICRO SIGN to ``\textmicro``.
2.5 (2023-09-12)
    - ``\textvarTheta`` is now an alias for ``\textTheta`` (the AMS-math
      command ``\varTheta`` sets the *letter* Theta in italic shape).
    - Fix errors in LuaTeX's 8-bit compatibility mode (luainputenc_).
    - Fix ``\MakeUppercase`` in PDF strings.
    - Drop composite definitions if the pre-composed character can also be
      selected by the `Unicode NFC normalization`_.
    - Test/fix case change commands with alphabeta_.
      Composite commands for PU.
      Inline ``alphabeta-tuenc.def`` and ``alphabeta-lgr.def``.
    - Update documentation, fix links.
2.6 (2023-11-16)
    -  Bugfix in alphabeta_: Don't use TextCommands for generic macros.


TODO:
    - Fix ``\textautosigma`` with Unicode fonts.

    .. report issues:
      The polytonic variant with dasia and oxia used in ἢ … ἤ (*either … or*)
      drops diacritics! By mistake, omission, or intent?

      Compilation error with MakeUppercase and combining ypogegrammeni in Greek
      locale: ``\foreignlanguage{greek}{Λͅ → \MakeUppercase{Λͅ}}``


__ http://tex.stackexchange.com/questions/281631/greek-small-rho-with-dasia-and-also-psili-problem-with-accent-and-lgr-encodin
.. _ltxbugs 4399:
   https://www.latex-project.org/cgi-bin/ltxbugs2html?pr=latex%2F4399&search=


References
==========

An alternative, more complete set of short mnemonic character names is
the `XML Entity Definitions for Characters`_ W3C Recommendation from
01 April 2010.

For glyph names of the LGR encoding see, e.g., ``CB.enc`` by Apostolos
Syropoulos and ``xl-lgr.enc`` from the libertine_ (legacy) package.
``lgr.cmap`` provides a mapping to Unicode characters.

A full set of ``\text*`` symbol macros is defined in ``ucsencs.def``
from the ucs_ package.

.. [babel-patch] Werner Lemberg, `Unicode support for the Greek LGR
   encoding` Εὔτυπον, τεῦχος  № 20, 2008.
   http://www.eutypon.gr/eutypon/pdf/e2008-20/e20-a03.pdf
.. [cbfonts] Claudio Beccari, `The CB Greek fonts`, Εὔτυπον, τεῦχος № 21, 2008.
   http://www.eutypon.gr/eutypon/pdf/e2008-21/e21-a01.pdf
.. [encguide] Frank Mittelbach, Robin Fairbairns, Werner Lemberg,
   LaTeX3 Project Team, `LaTeX font encodings`.
   https://mirrors.ctan.org/macros/latex/base/encguide.pdf
.. [fntguide] LaTeX3 Project Team, `LaTeX2ε font selection`.
   https://mirrors.ctan.org/macros/latex/base/fntguide.pdf
.. [greek-usage] Apostolos Syropoulos, `Writing Greek with the greek option
   of the babel package`, 1997.
   https://mirrors.ctan.org/language/babel/contrib/greek/usage.pdf
.. [ltnews26] LaTeX Project Team, `LaTeX News` Issue 26, January 2017.
   https://www.latex-project.org/news/latex2e-news/ltnews26.pdf
.. [ltnews31] `LaATeX News`, Issue 31, February 2020, p. 3:
   https://www.latex-project.org/news/latex2e-news/ltnews31.pdf.
.. [teubner-doc] Claudio Beccari, ``teubner.sty``
   `An extension to the greek option of the babel package`, 2011.
   https://mirrors.ctan.org/macros/latex/contrib/teubner/teubner-doc.pdf

.. _LaTeX Project Public License: http://www.latex-project.org/lppl.txt
.. _PyLit: https://pypi.org/project/pylit/
.. _reStructuredText: https://docutils.sourceforge.io/rst.html
.. _Docutils: https://docutils.sourceforge.io/rst.html

.. _Adobe Glyph List For New Fonts:
    http://raw.githubusercontent.com/adobe-type-tools/agl-aglfn/master/aglfn.txt
.. _Adobe Glyph List:
    http://raw.githubusercontent.com/adobe-type-tools/agl-aglfn/master/glyphlist.txt
.. _Unicode names list: http://www.unicode.org/Public/UNIDATA/NamesList.txt
.. _Unicode NFC normalization: https://www.unicode.org/reports/tr15/
.. _XML Entity Definitions for Characters:
    http://www.w3.org/TR/xml-entity-names/
.. _CB fonts: https://ctan.org/pkg/cbgreek-complete
.. _CM Unicode: https://ctan.org/pkg/cm-unicode
.. _Deja Vu: http://dejavu-fonts.org
.. _EB Garamond: https://ctan.org/pkg/ebgaramond
.. _GFS fonts: https://ctan.org/pkg/gfs
.. _Kerkis: https://ctan.org/pkg/kerkis
.. _Latin Modern: http://www.gust.org.pl/projects/e-foundry/latin-modern
.. _Libertine OTF: https://ctan.org/pkg/libertineotf
.. _Libertine GC: https://ctan.org/pkg/libertinegc
.. _Libertinus: https://ctan.org/pkg/libertinus
.. _Old Standard: https://ctan.org/pkg/oldstandard
.. _OpenType Font: https://ctan.org/topic/font-otf
.. _Tempora: https://ctan.org/pkg/tempora
.. _UM Typewriter: https://ctan.org/pkg/umtypewriter
.. _amssymb: https://ctan.org/pkg/amsfonts
.. _arabi: https://ctan.org/pkg/arabi
.. _babel-greek: https://ctan.org/pkg/babel-greek
.. _babel: https://ctan.org/pkg/babel
.. _dejavu: https://ctan.org/pkg/dejavu
.. _fontenc:  https://ctan.org/pkg/fontenc
.. _fontspec:  https://ctan.org/pkg/fontspec
.. _greek-inputenc: https://ctan.org/pkg/greek-inputenc
.. _grtimes: https://ctan.org/pkg/grtimes
.. _hyperref: https://ctan.org/pkg/hyperref
.. _ibycus-babel: https://ctan.org/pkg/ibycus-babel
.. _ibygrk: https://ctan.org/pkg/ibygrk
.. _lgrx: https://ctan.org/pkg/lgrx
.. _libertine: https://ctan.org/pkg/libertine-legacy
.. _luainputenc: https://ctan.org/pkg/luainputenc
.. _substitutefont: https://ctan.org/pkg/substitutefont
.. _teubner: https://ctan.org/pkg/teubner
.. _ucs: https://ctan.org/pkg/unicode
.. _unicode-math: https://ctan.org/pkg/unicode-math
.. _xgreek: https://ctan.org/pkg/xgreek
.. _xunicode: https://ctan.org/pkg/xunicode
