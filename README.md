# Graphgenerator
Appendix to my Bachelor's thesis

This repository contains the digital appendix for my bachelor's thesis "Lower Bounds for the Weisfeiler-Leman Dimension of Graphs of Bounded Tree Width".

It consists of:
- An implementation: This is a Go Script for generating grids, walls and rectangular grids as well as CFI graphs over these instances and, for experimental purposes, pyramid graphs.
- The test instances: Walls, Rectangles, CFI Grids, CFI Rectangles and CFI Walls up to dimension 20 as well as Grids up to dimension 100
- Results: These are the results of the tested instances (which can be found in the .td files) and the output files of the RWTH Cluster.

The graphs are generated in a PACE2017 compatible .gr format to be used with the PID-BT algorithm by Hisao Tamaki.
