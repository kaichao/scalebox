# Configuration file for the Sphinx documentation builder.
#
# For the full list of built-in configuration values, see the documentation:
# https://www.sphinx-doc.org/en/master/usage/configuration.html

# -- Project information -----------------------------------------------------
# https://www.sphinx-doc.org/en/master/usage/configuration.html#project-information

project = 'scalebox'
copyright = '2024, scalebox'
author = 'scalebox'
release = '0.1.0'

# -- General configuration ---------------------------------------------------
# https://www.sphinx-doc.org/en/master/usage/configuration.html#general-configuration

extensions = [
    'recommonmark',
    'myst-parser'
    'sphinx_markdown_tables',
    'sphinxcontrib.mermaid',
    'sphinx.ext.autosectionlabel'
]

templates_path = ['_templates']
exclude_patterns = ['_build', 'Thumbs.db', '.DS_Store']

language = 'zh_CN'

# -- Options for HTML output -------------------------------------------------
# https://www.sphinx-doc.org/en/master/usage/configuration.html#options-for-html-output

html_theme = 'sphinx_rtd_theme'
html_static_path = ['_static']

# html_build_dir = os.environ.get('READTHEDOCS_OUTPUT', 'docs/en/build/html')

#supported markdown
source_suffix = {
    '.rst': 'restructuredtext',
    '.md': 'markdown'
}

master_doc = 'index'
