.. Scalebox documentation master file, created by
   sphinx-quickstart on Sun Jan 21 09:48:07 2024.
   You can adapt this file completely to your liking, but it should at least
   contain the root `toctree` directive.

Scalebox 中文文档
====================================

本中文站面向用户、应用开发者与系统管理员；AI agent 工作流与机器可读参考见仓库 `docs/agents/`，应用仓接入方法见 :doc:`AI 辅助开发 <application/9_ai_assisted_development>`。

根据您的使用场景，选择以下学习路径：

- **应用开发者**：:doc:`入门教程 <started/1_introduction>` → :doc:`核心概念 <application/1_core_concepts>` → :doc:`运行与管理应用 <application/2_running_apps>` → :doc:`应用设计 <application/3_app_design>` → :doc:`模块开发 <application/4_module_development>`
- **部署与运维**：:doc:`安装部署 <deployment/1_installation>` → :doc:`运维管理 <deployment/2_operations>` → :doc:`FAQ <appendix/3_faq>`
- **平台贡献者 / 深入原理**：:doc:`编程模型 <design/1_programming_model>` → :doc:`主路由与应用状态 <design/2_main_router_status>` → :doc:`VTask 虚拟任务 <design/3_vtask>` → :doc:`节点本地计算 <design/4_node_local_compute>` → :doc:`跨集群架构 <design/5_cross_cluster>` → :doc:`安全框架 <design/6_security>` → :doc:`高级编程特性 <design/7_advanced_features>`

**开发者与 AI 参考（仓库直读）**

中文站面向人（用户、应用开发者、系统管理员）；AI agent 工作流与机器可读参考在仓库中直读：

- `AGENTS.md <https://github.com/kaichao/scalebox/blob/dev/AGENTS.md>`_
- `docs/agents/ <https://github.com/kaichao/scalebox/tree/dev/docs/agents>`_
- `features/ <https://github.com/kaichao/scalebox/tree/dev/features>`_
- `docs/reference/ <https://github.com/kaichao/scalebox/tree/dev/docs/reference>`_

.. toctree::
   :maxdepth: 2
   :caption: 一、入门教程
   :hidden:

   started/1_introduction
   started/2_quick_start
   started/3_first_app

.. toctree::
   :maxdepth: 2
   :caption: 二、应用开发
   :hidden:

   application/1_core_concepts
   application/2_running_apps
   application/3_app_design
   application/4_module_development
   application/5_example_apps
   application/6_standard_modules
   application/7_webui
   application/8_vscode
   application/9_ai_assisted_development

.. toctree::
   :maxdepth: 2
   :caption: 三、部署与运维
   :hidden:

   deployment/1_installation
   deployment/2_operations

.. toctree::
   :maxdepth: 2
   :caption: 四、设计与原理
   :hidden:

   design/1_programming_model
   design/2_main_router_status
   design/3_vtask
   design/4_node_local_compute
   design/5_cross_cluster
   design/6_security
   design/7_advanced_features

.. toctree::
   :maxdepth: 2
   :caption: 五、参考
   :hidden:

   reference/1_tech_specifications
   reference/2_parameter_reference
   reference/3_commandline_tools
   reference/4_exit_code_spec

.. toctree::
   :maxdepth: 2
   :caption: 六、附录与其他
   :hidden:

   appendix/1_shell_programming
   appendix/2_software_install
   appendix/3_faq
   appendix/4_contributors
   appendix/5_contact_us

Indices and tables
==================

* :ref:`genindex`
* :ref:`modindex`
* :ref:`search`
